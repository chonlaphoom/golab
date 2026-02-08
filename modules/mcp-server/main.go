package main

import (
	"bufio"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	protocolVersion = "2025-03-26" // latest MCP protocol version supported
	maxWorkers      = 5            // TODO: add timeout for each worker task
)

type ThreadSafeWriter struct {
	writer *bufio.Writer
	mu     sync.Mutex
}

func main() {
	setupLog()
	log.Printf("MCP Server %s is running...", protocolVersion)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wg := &sync.WaitGroup{}

	decoder := json.NewDecoder(os.Stdin)
	safeWriter := &ThreadSafeWriter{writer: bufio.NewWriter(os.Stdout)}

	msgChan := make(chan json.RawMessage)
	errChan := make(chan error)

	go readAndPushMsgs(ctx, cancel, decoder, msgChan, errChan)

	for i := range maxWorkers {
		wg.Add(1)
		go worker(ctx, wg, safeWriter, msgChan, errChan, i)
	}

	wg.Wait()
	log.Println("All workers have exited, shutting down MCP server.")
}
