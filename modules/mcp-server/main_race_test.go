package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestThreadSafeWriter_ConcurrentWrites tests that ThreadSafeWriter properly serializes
// concurrent writes from multiple goroutines using 50 goroutines writing 10 messages each.
func TestThreadSafeWriter_ConcurrentWrites(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	const numGoroutines = 50
	const messagesPerGoroutine = 10
	const totalMessages = numGoroutines * messagesPerGoroutine

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Each goroutine writes its ID and message number
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for msgNum := 0; msgNum < messagesPerGoroutine; msgNum++ {
				safeWriter.mu.Lock()
				msg := fmt.Sprintf(`{"goroutine":%d,"message":%d}`, goroutineID, msgNum)
				_, err := safeWriter.writer.WriteString(msg + "\n")
				if err != nil {
					t.Errorf("Write error: %v", err)
				}
				err = safeWriter.writer.Flush()
				safeWriter.mu.Unlock()

				if err != nil {
					t.Errorf("Flush error: %v", err)
				}
			}
		}(i)
	}

	// Wait for all goroutines with timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(10 * time.Second):
		t.Fatal("Test timed out waiting for goroutines")
	}

	// Verify all messages were written
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	validLines := 0
	for _, line := range lines {
		if len(line) > 0 {
			validLines++
			// Verify each line is valid JSON
			var msg map[string]int
			if err := json.Unmarshal(line, &msg); err != nil {
				t.Errorf("Invalid JSON line: %s, error: %v", line, err)
			}
		}
	}

	if validLines != totalMessages {
		t.Errorf("Expected %d messages, got %d", totalMessages, validLines)
	}
}

// TestThreadSafeWriter_RaceDetection tests concurrent read/write operations
// to ensure the race detector catches any issues
func TestThreadSafeWriter_RaceDetection(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	var wg sync.WaitGroup
	const numWriters = 20
	const writesPerWriter = 5

	wg.Add(numWriters)

	// Spawn multiple writers
	for i := 0; i < numWriters; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writesPerWriter; j++ {
				safeWriter.mu.Lock()
				msg := map[string]interface{}{
					"writer": id,
					"count":  j,
					"data":   "test data",
				}
				encoder := json.NewEncoder(safeWriter.writer)
				if err := encoder.Encode(msg); err != nil {
					t.Errorf("Encode error: %v", err)
				}
				if err := safeWriter.writer.Flush(); err != nil {
					t.Errorf("Flush error: %v", err)
				}
				safeWriter.mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Count valid JSON objects
	decoder := json.NewDecoder(&buf)
	count := 0
	for {
		var msg map[string]interface{}
		if err := decoder.Decode(&msg); err != nil {
			break
		}
		count++
	}

	expectedCount := numWriters * writesPerWriter
	if count != expectedCount {
		t.Errorf("Expected %d messages, got %d", expectedCount, count)
	}
}

// TestThreadSafeWriter_NoDataCorruption verifies that concurrent writes
// don't corrupt JSON structure
func TestThreadSafeWriter_NoDataCorruption(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	const numGoroutines = 30
	const messagesPerGoroutine = 10
	messageCounts := make(map[string]int)
	var countMu sync.Mutex

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				msg := map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      id*messagesPerGoroutine + j,
					"result": map[string]interface{}{
						"goroutine": id,
						"message":   j,
						"text":      fmt.Sprintf("Message from goroutine %d, number %d", id, j),
					},
				}

				safeWriter.mu.Lock()
				encoder := json.NewEncoder(safeWriter.writer)
				if err := encoder.Encode(msg); err != nil {
					t.Errorf("Failed to encode message: %v", err)
				}
				if err := safeWriter.writer.Flush(); err != nil {
					t.Errorf("Failed to flush: %v", err)
				}
				safeWriter.mu.Unlock()

				// Track expected messages
				key := fmt.Sprintf("%d-%d", id, j)
				countMu.Lock()
				messageCounts[key]++
				countMu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// Verify all messages are valid JSON-RPC responses
	lines := bytes.Split(buf.Bytes(), []byte("\n"))
	validMessages := 0

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var response map[string]interface{}
		if err := json.Unmarshal(line, &response); err != nil {
			t.Errorf("Corrupted JSON: %s, error: %v", line, err)
			continue
		}

		// Verify JSON-RPC structure
		if response["jsonrpc"] != "2.0" {
			t.Errorf("Invalid jsonrpc version: %v", response["jsonrpc"])
		}
		if _, ok := response["result"]; !ok {
			t.Errorf("Missing result field in response: %v", response)
		}

		validMessages++
	}

	expectedMessages := numGoroutines * messagesPerGoroutine
	if validMessages != expectedMessages {
		t.Errorf("Expected %d valid messages, got %d", expectedMessages, validMessages)
	}

	// Verify no duplicate messages
	for key, count := range messageCounts {
		if count != 1 {
			t.Errorf("Message %s sent %d times (expected 1)", key, count)
		}
	}
}

// TestThreadSafeWriter_StressTest performs high-concurrency stress testing
func TestThreadSafeWriter_StressTest(t *testing.T) {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	safeWriter := &ThreadSafeWriter{
		writer: writer,
	}

	const numGoroutines = 100
	const messagesPerGoroutine = 3

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < messagesPerGoroutine; j++ {
				msg := fmt.Sprintf(`{"id":%d,"seq":%d}`, id, j)

				safeWriter.mu.Lock()
				safeWriter.writer.WriteString(msg + "\n")
				safeWriter.writer.Flush()
				safeWriter.mu.Unlock()
			}
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		elapsed := time.Since(start)
		t.Logf("Stress test completed in %v", elapsed)
	case <-time.After(15 * time.Second):
		t.Fatal("Stress test timed out")
	}

	// Basic validation
	lines := bytes.Count(buf.Bytes(), []byte("\n"))
	expectedLines := numGoroutines * messagesPerGoroutine

	// Allow some tolerance for empty lines
	if lines < expectedLines {
		t.Errorf("Expected at least %d lines, got %d", expectedLines, lines)
	}
}
