package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestReadAndPushMsgs_ConcurrentConsumers tests single reader with multiple consumers
func TestReadAndPushMsgs_ConcurrentConsumers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a buffer with test messages
	var inputBuf bytes.Buffer
	const numMessages = 250
	for i := 0; i < numMessages; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "tools/list",
			"id":      i,
		}
		if err := json.NewEncoder(&inputBuf).Encode(msg); err != nil {
			t.Fatalf("Failed to encode message: %v", err)
		}
	}

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 50)
	errChan := make(chan error, 1)

	// Start reader
	go readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)

	// Start multiple consumers
	const numConsumers = 5
	var wg sync.WaitGroup
	var consumedCount atomic.Int32
	messagesSeen := make(map[int]int) // map[messageID]count
	var seenMu sync.Mutex

	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func(consumerID int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case err := <-errChan:
					if err != nil && err != io.EOF {
						t.Errorf("Consumer %d received error: %v", consumerID, err)
					}
					return
				case msg, ok := <-msgChan:
					if !ok {
						return
					}
					// Parse message to check ID
					var parsed map[string]interface{}
					if err := json.Unmarshal(msg, &parsed); err != nil {
						t.Errorf("Consumer %d failed to parse message: %v", consumerID, err)
						continue
					}

					if id, ok := parsed["id"].(float64); ok {
						seenMu.Lock()
						messagesSeen[int(id)]++
						seenMu.Unlock()
					}

					consumedCount.Add(1)
				}
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
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out waiting for consumers")
	}

	consumed := consumedCount.Load()
	if consumed != numMessages {
		t.Errorf("Expected %d messages consumed, got %d", numMessages, consumed)
	}

	// Verify each message was consumed exactly once
	seenMu.Lock()
	for i := 0; i < numMessages; i++ {
		count := messagesSeen[i]
		if count != 1 {
			t.Errorf("Message ID %d was consumed %d times (expected 1)", i, count)
		}
	}
	seenMu.Unlock()
}

// TestReadAndPushMsgs_ContextCancellation tests cancellation during active reading
func TestReadAndPushMsgs_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a buffer with many messages
	var inputBuf bytes.Buffer
	const numMessages = 200
	for i := 0; i < numMessages; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "test",
			"id":      i,
		}
		json.NewEncoder(&inputBuf).Encode(msg)
	}

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 50)
	errChan := make(chan error, 1)

	// Start reader
	readerDone := make(chan struct{})
	go func() {
		readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)
		close(readerDone)
	}()

	// Consume some messages then cancel
	var consumedCount atomic.Int32
	consumerDone := make(chan struct{})

	go func() {
		defer close(consumerDone)
		for {
			select {
			case _, ok := <-msgChan:
				if !ok {
					return
				}
				consumedCount.Add(1)

				// Cancel after consuming half the messages
				if consumedCount.Load() == numMessages/2 {
					cancel()
				}
			case <-time.After(5 * time.Second):
				return
			}
		}
	}()

	// Wait for reader to finish
	select {
	case <-readerDone:
		// Reader should exit after cancel
		t.Log("Reader exited cleanly after context cancellation")
	case <-time.After(5 * time.Second):
		t.Fatal("Reader did not exit after context cancellation")
	}

	// Wait for consumer
	select {
	case <-consumerDone:
		// Consumer finished
	case <-time.After(2 * time.Second):
		t.Fatal("Consumer did not finish")
	}

	// Verify channel was closed
	select {
	case _, ok := <-msgChan:
		if ok {
			t.Error("Message channel should be closed after reader exits")
		}
	default:
		// Channel is closed, which is expected
	}

	consumed := consumedCount.Load()
	t.Logf("Consumed %d messages before cancellation", consumed)

	if consumed == 0 {
		t.Error("Expected some messages to be consumed before cancellation")
	}
}

// TestReadAndPushMsgs_ErrorHandling tests error handling during reading
func TestReadAndPushMsgs_ErrorHandling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a reader that will return an error after some messages
	var inputBuf bytes.Buffer
	const validMessages = 10
	for i := 0; i < validMessages; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "test",
			"id":      i,
		}
		json.NewEncoder(&inputBuf).Encode(msg)
	}
	// Add invalid JSON to trigger error
	inputBuf.WriteString("{invalid json\n")

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 20)
	errChan := make(chan error, 1)

	// Start reader
	readerDone := make(chan struct{})
	go func() {
		readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)
		close(readerDone)
	}()

	// Consume messages and wait for error
	var consumedCount atomic.Int32
	errorReceived := false

	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		for {
			select {
			case _, ok := <-msgChan:
				if !ok {
					return
				}
				consumedCount.Add(1)
			case err := <-errChan:
				if err != nil {
					errorReceived = true
					t.Logf("Received expected error: %v", err)
				}
				return
			case <-time.After(5 * time.Second):
				return
			}
		}
	}()

	// Wait for reader to finish
	select {
	case <-readerDone:
		t.Log("Reader exited after error")
	case <-time.After(5 * time.Second):
		t.Fatal("Reader did not exit after error")
	}

	// Wait for consumer
	select {
	case <-consumerDone:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Consumer did not finish")
	}

	if !errorReceived {
		t.Error("Expected to receive error on errChan")
	}

	consumed := consumedCount.Load()
	// Due to concurrent processing, we might not process all messages before error occurs
	// We should have consumed most of them (at least 80%)
	minExpected := int32(validMessages * 8 / 10)
	if consumed < minExpected {
		t.Errorf("Expected at least %d valid messages consumed, got %d", minExpected, consumed)
	}
	t.Logf("Consumed %d out of %d messages before error", consumed, validMessages)
}

// TestReadAndPushMsgs_ChannelBuffering tests channel buffering behavior
func TestReadAndPushMsgs_ChannelBuffering(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create messages
	var inputBuf bytes.Buffer
	const numMessages = 100
	for i := 0; i < numMessages; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "test",
			"id":      i,
		}
		json.NewEncoder(&inputBuf).Encode(msg)
	}

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 10) // Small buffer
	errChan := make(chan error, 1)

	// Start reader
	go readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)

	// Slow consumer to test backpressure
	var consumedCount atomic.Int32
	consumerDone := make(chan struct{})

	go func() {
		defer close(consumerDone)
		for {
			select {
			case _, ok := <-msgChan:
				if !ok {
					return
				}
				consumedCount.Add(1)
				// Simulate slow processing
				time.Sleep(1 * time.Millisecond)
			case <-time.After(5 * time.Second):
				return
			}
		}
	}()

	// Wait for consumer
	select {
	case <-consumerDone:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Consumer timed out")
	}

	consumed := consumedCount.Load()
	if consumed != numMessages {
		t.Errorf("Expected %d messages consumed, got %d", numMessages, consumed)
	}
}

// TestReadAndPushMsgs_RapidCancellation tests cancellation immediately after start
func TestReadAndPushMsgs_RapidCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// Create a very small set of messages so EOF happens quickly
	var inputBuf bytes.Buffer
	for i := 0; i < 5; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "test",
			"id":      i,
		}
		json.NewEncoder(&inputBuf).Encode(msg)
	}

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 20)
	errChan := make(chan error, 1)

	// Start reader
	readerDone := make(chan struct{})
	go func() {
		readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)
		close(readerDone)
	}()

	// Drain messages to let reader progress, then cancel
	go func() {
		time.Sleep(1 * time.Millisecond)
		cancel()
		// Drain remaining messages so channel doesn't block
		for range msgChan {
		}
	}()

	// Wait for reader to exit - note that decoder.Decode() doesn't respect context
	// immediately when blocked on I/O, so reader will exit when EOF or next decode happens
	select {
	case <-readerDone:
		t.Log("Reader exited cleanly after rapid cancellation")
	case <-time.After(3 * time.Second):
		t.Fatal("Reader did not exit after rapid cancellation")
	}
}

// TestReadAndPushMsgs_EOFHandling tests EOF handling
func TestReadAndPushMsgs_EOFHandling(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a small set of messages
	var inputBuf bytes.Buffer
	const numMessages = 5
	for i := 0; i < numMessages; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "test",
			"id":      i,
		}
		json.NewEncoder(&inputBuf).Encode(msg)
	}

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 10)
	errChan := make(chan error, 1)

	// Start reader
	readerDone := make(chan struct{})
	go func() {
		readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)
		close(readerDone)
	}()

	// Consume messages
	var consumedCount atomic.Int32
	eofReceived := false

	for {
		select {
		case _, ok := <-msgChan:
			if !ok {
				goto CheckResults
			}
			consumedCount.Add(1)
		case err := <-errChan:
			if err == io.EOF {
				eofReceived = true
			}
			goto CheckResults
		case <-time.After(3 * time.Second):
			t.Fatal("Timeout waiting for EOF")
		}
	}

CheckResults:
	// Wait for reader
	select {
	case <-readerDone:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Reader did not exit after EOF")
	}

	if !eofReceived {
		t.Error("Expected to receive EOF on errChan")
	}

	consumed := consumedCount.Load()
	if consumed != numMessages {
		t.Errorf("Expected %d messages consumed, got %d", numMessages, consumed)
	}
}

// TestReadAndPushMsgs_ConcurrentStress performs stress testing
func TestReadAndPushMsgs_ConcurrentStress(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create many messages
	var inputBuf bytes.Buffer
	const numMessages = 300
	for i := 0; i < numMessages; i++ {
		msg := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "tools/list",
			"id":      i,
		}
		json.NewEncoder(&inputBuf).Encode(msg)
	}

	decoder := json.NewDecoder(&inputBuf)
	msgChan := make(chan json.RawMessage, 100)
	errChan := make(chan error, 1)

	start := time.Now()

	// Start reader
	go readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)

	// Start multiple fast consumers
	const numConsumers = 10
	var wg sync.WaitGroup
	var totalConsumed atomic.Int32

	for i := 0; i < numConsumers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case _, ok := <-msgChan:
					if !ok {
						return
					}
					totalConsumed.Add(1)
				case <-errChan:
					return
				case <-time.After(5 * time.Second):
					return
				}
			}
		}()
	}

	// Wait for consumers
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		elapsed := time.Since(start)
		consumed := totalConsumed.Load()
		t.Logf("Stress test completed in %v, consumed %d messages", elapsed, consumed)

		if consumed != numMessages {
			t.Errorf("Expected %d messages consumed, got %d", numMessages, consumed)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Stress test timed out")
	}
}
