package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestWorker_ConcurrentMessageProcessing tests that multiple workers can
// process messages concurrently without race conditions
func TestWorker_ConcurrentMessageProcessing(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	msgChan := make(chan json.RawMessage, 100)
	errChan := make(chan error, 1)

	const numWorkers = 5
	const numMessages = 300

	var wg sync.WaitGroup
	var processedCount atomic.Int32

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					// Drain remaining messages
					for range msgChan {
						processedCount.Add(1)
					}
					return
				case err := <-errChan:
					if err != nil {
						t.Errorf("Worker %d received error: %v", workerID, err)
					}
					return
				case msg, ok := <-msgChan:
					if !ok {
						return
					}
					// Process message (use safeWriter)
					processMessage(safeWriter, msg)
					processedCount.Add(1)
				}
			}
		}(i)
	}

	// Send messages
	go func() {
		for i := 0; i < numMessages; i++ {
			msg := json.RawMessage([]byte(`{"jsonrpc":"2.0","method":"tools/list","id":` + string(rune(i)) + `}`))
			msgChan <- msg
		}
		close(msgChan)
	}()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out waiting for workers")
	}

	processed := processedCount.Load()
	if processed != numMessages {
		t.Errorf("Expected %d messages processed, got %d", numMessages, processed)
	}
}

// TestWorker_ErrorChannelRace tests that error channel handling doesn't cause races
func TestWorker_ErrorChannelRace(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	msgChan := make(chan json.RawMessage, 50)
	errChan := make(chan error, 10) // Buffered to prevent blocking

	const numWorkers = 5

	var wg sync.WaitGroup
	var errorCount atomic.Int32

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case err, ok := <-errChan:
					if !ok {
						return
					}
					if err != nil {
						errorCount.Add(1)
					}
					return
				case msg, ok := <-msgChan:
					if !ok {
						return
					}
					processMessage(safeWriter, msg)
				}
			}
		}(i)
	}

	// Send some messages and then an error
	go func() {
		for i := 0; i < 20; i++ {
			msgChan <- json.RawMessage(`{"jsonrpc":"2.0","method":"test"}`)
		}
		// Send error to trigger shutdown
		errChan <- nil // EOF simulation
		close(msgChan)
	}()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success - workers shut down cleanly
		t.Logf("Workers shut down cleanly, errors received: %d", errorCount.Load())
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out waiting for workers to handle error")
	}
}

// TestWorker_GracefulShutdown tests that workers drain messages during shutdown
func TestWorker_GracefulShutdown(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	msgChan := make(chan json.RawMessage, 200)
	errChan := make(chan error, 1)

	const numWorkers = 5
	const numMessages = 200

	var wg sync.WaitGroup
	var processedBeforeCancel atomic.Int32
	var processedAfterCancel atomic.Int32

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			cancelReceived := false

			for {
				select {
				case <-ctx.Done():
					cancelReceived = true
					// Drain phase: process remaining messages
					for msg := range msgChan {
						processMessage(safeWriter, msg)
						processedAfterCancel.Add(1)
					}
					return
				case err := <-errChan:
					if err != nil {
						t.Errorf("Worker %d error: %v", workerID, err)
					}
					return
				case msg, ok := <-msgChan:
					if !ok {
						return
					}
					processMessage(safeWriter, msg)
					if !cancelReceived {
						processedBeforeCancel.Add(1)
					} else {
						processedAfterCancel.Add(1)
					}
				}
			}
		}(i)
	}

	// Send messages
	go func() {
		for i := 0; i < numMessages; i++ {
			msg := json.RawMessage(`{"jsonrpc":"2.0","method":"tools/list","id":` + string(rune(i)) + `}`)
			msgChan <- msg

			// Cancel halfway through
			if i == numMessages/2 {
				time.Sleep(10 * time.Millisecond) // Let some messages process
				cancel()
			}
		}
		close(msgChan)
	}()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out during graceful shutdown")
	}

	beforeCancel := processedBeforeCancel.Load()
	afterCancel := processedAfterCancel.Load()
	total := beforeCancel + afterCancel

	t.Logf("Processed before cancel: %d, after cancel: %d, total: %d", beforeCancel, afterCancel, total)

	if total != numMessages {
		t.Errorf("Expected %d total messages processed, got %d", numMessages, total)
	}

	if afterCancel == 0 {
		t.Error("Expected some messages to be drained after cancel, got 0")
	}
}

// TestWorker_ChannelClosingRace tests race between channel closing and worker reads
func TestWorker_ChannelClosingRace(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	msgChan := make(chan json.RawMessage, 100)
	errChan := make(chan error, 1)

	const numWorkers = 5
	const numMessages = 100

	var wg sync.WaitGroup
	var processedCount atomic.Int32
	var cleanExits atomic.Int32

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			defer cleanExits.Add(1)

			for {
				select {
				case <-ctx.Done():
					for msg := range msgChan {
						processMessage(safeWriter, msg)
						processedCount.Add(1)
					}
					return
				case err := <-errChan:
					if err != nil {
						t.Errorf("Worker %d error: %v", workerID, err)
					}
					return
				case msg, ok := <-msgChan:
					if !ok {
						// Channel closed, exit cleanly
						return
					}
					processMessage(safeWriter, msg)
					processedCount.Add(1)
				}
			}
		}(i)
	}

	// Send messages and close channel while workers are reading
	go func() {
		for i := 0; i < numMessages; i++ {
			msg := json.RawMessage(`{"jsonrpc":"2.0","method":"tools/list","id":` + string(rune(i)) + `}`)
			select {
			case msgChan <- msg:
			case <-time.After(1 * time.Second):
				t.Error("Timeout sending message")
				return
			}

			// Close channel partway through to test race condition
			if i == numMessages/2 {
				time.Sleep(5 * time.Millisecond)
				close(msgChan)
				return
			}
		}
	}()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out during channel closing race test")
	}

	processed := processedCount.Load()
	exits := cleanExits.Load()

	t.Logf("Messages processed: %d, clean exits: %d", processed, exits)

	if exits != numWorkers {
		t.Errorf("Expected %d workers to exit cleanly, got %d", numWorkers, exits)
	}

	// We should have processed at least the messages sent before close
	expectedMin := int32(numMessages / 2)
	if processed < expectedMin {
		t.Errorf("Expected at least %d messages processed, got %d", expectedMin, processed)
	}
}

// TestWorker_ProcessMessageConcurrency tests processMessage function under concurrent load
func TestWorker_ProcessMessageConcurrency(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	const numGoroutines = 50
	const messagesPerGoroutine = 5

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Valid JSON-RPC messages
	messages := []json.RawMessage{
		json.RawMessage(`{"jsonrpc":"2.0","method":"initialize","params":{"protocolVersion":"2025-03-26"},"id":1}`),
		json.RawMessage(`{"jsonrpc":"2.0","method":"tools/list","id":2}`),
		json.RawMessage(`{"jsonrpc":"2.0","method":"tools/call","params":{"name":"echo","arguments":{"text":"test"}},"id":3}`),
		json.RawMessage(`{"jsonrpc":"2.0","method":"notifications/initialized"}`),
		json.RawMessage(`{"jsonrpc":"2.0","method":"tools/list","id":4}`),
	}

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				msg := messages[j%len(messages)]
				processMessage(safeWriter, msg)
			}
		}(i)
	}

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
		t.Logf("All goroutines completed successfully")
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out during processMessage concurrency test")
	}

	// Verify output contains valid JSON responses
	// The output should contain responses, though order is not guaranteed
	if buf.Len() == 0 {
		t.Error("Expected some output from processMessage calls")
	}
}

// TestWorker_HighLoadStress performs stress testing with high message volume
func TestWorker_HighLoadStress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	msgChan := make(chan json.RawMessage, 500)
	errChan := make(chan error, 1)

	const numWorkers = 5
	const numMessages = 500

	var wg sync.WaitGroup
	var processedCount atomic.Int32

	start := time.Now()

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					for msg := range msgChan {
						processMessage(safeWriter, msg)
						processedCount.Add(1)
					}
					return
				case err := <-errChan:
					if err != nil {
						return
					}
					return
				case msg, ok := <-msgChan:
					if !ok {
						return
					}
					processMessage(safeWriter, msg)
					processedCount.Add(1)
				}
			}
		}(i)
	}

	// Send messages rapidly
	go func() {
		for i := 0; i < numMessages; i++ {
			msg := json.RawMessage(`{"jsonrpc":"2.0","method":"tools/list","id":` + string(rune(i)) + `}`)
			msgChan <- msg
		}
		close(msgChan)
	}()

	// Wait with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		elapsed := time.Since(start)
		processed := processedCount.Load()
		t.Logf("High load stress test completed in %v, processed %d messages", elapsed, processed)

		if processed != numMessages {
			t.Errorf("Expected %d messages processed, got %d", numMessages, processed)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("Stress test timed out")
	}
}
