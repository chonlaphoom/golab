package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mcp-server/jsonrpc2"
	"mcp-server/mcp"
	"os"
	"os/signal"
	"syscall"
)

const protocolVersion = "2026-01-26"

func main() {
	log.Printf("MCP Server %s is running...", protocolVersion)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	decoder := json.NewDecoder(os.Stdin)
	// writer := bufio.NewWriter(os.Stdout)

	msgChan := make(chan json.RawMessage)
	errChan := make(chan error)

	go func() {
		for {
			var msg json.RawMessage
			if err := decoder.Decode(&msg); err != nil {
				errChan <- err
				return
			}
			msgChan <- msg
		}
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("Received shutdown signal, exiting.")
			return
		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				log.Println("EOF received, exiting.")
				break
			}
			log.Fatalf("Error reading message: %v", err)
		case msg := <-msgChan:
			log.Printf("Received message: %s", string(msg))
			p := jsonrpc2.NewParser()
			if err := p.ParseRequest(msg); err != nil {
				log.Printf("Error parsing request: %v, wait for next message", err)
				continue
			}

			var err error
			switch p.Req.Method {
			case "initialize":
				mcp.HandleInitialize(protocolVersion)
			case "tools/list":
				mcp.HandleListTools()
			case "tool/call":
			}
			if err != nil {
				log.Printf("Error handling method %s: %v", p.Req.Method, err)
				continue
			}
			// writer.Write(responseBytes)
			// writer.Flush()
			// log.Printf("Sent response: %s", string(responseBytes))
		}
	}
}
