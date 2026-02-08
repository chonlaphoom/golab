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

func main() {
	setupLog()
	log.Printf("MCP Server %s is running...", protocolVersion)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	wg := &sync.WaitGroup{}

	decoder := json.NewDecoder(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	msgChan := make(chan json.RawMessage)
	errChan := make(chan error)

	go readAndPushMsgs(ctx, decoder, msgChan, errChan)

	for range maxWorkers {
		wg.Add(1)
		go worker(ctx, wg, writer, msgChan, errChan)
	}

	wg.Wait()
	log.Println("All workers have exited, shutting down MCP server.")
}
