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

const protocolVersion = "2025-03-26" // latest MCP protocol version supported

func main() {
	log.Printf("MCP Server %s is running...", protocolVersion)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	decoder := json.NewDecoder(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)

	msgChan := make(chan json.RawMessage)
	errChan := make(chan error)

	go func() {
		defer close(msgChan)
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
				return
			}
			log.Fatalf("Error reading message: %v", err)
		case msg, ok := <-msgChan:
			if !ok {
				log.Println("Message channel closed, exiting.")
				return
			}

			log.Printf("Message received on channel: %s", string(msg))

			log.Printf("Received message: %s", string(msg))
			p := jsonrpc2.NewParser()
			if err := p.ParseRequest(msg); err != nil {
				log.Printf("Error parsing request: %v, wait for next message", err)
				continue
			}

			var err error
			var res any
			switch p.Req.Method {
			case "initialize":
				v, ok := p.Req.Params.GetAsObject()
				if ok {
					log.Println("Initialize params:", v)
				}
				res = mcp.HandleInitialize(protocolVersion, p.Req.ID)
			case "tools/list":
				res = mcp.HandleListTools(p.Req.ID)
			case "tools/call":
				res = mcp.HandleCallTool(p.Req.Params, p.Req.ID)
			default:
				err = errors.New("unknown method: " + p.Req.Method)
			}
			if err != nil {
				log.Printf("Error handling method %s: %v", p.Req.Method, err)
				continue
			}

			en := json.NewEncoder(writer)
			if err := en.Encode(res); err != nil {
				log.Fatalf("Error encoding response: %v", err)
				continue
			}

			if err := writer.Flush(); err != nil {
				log.Fatalf("Error flushing writer: %v", err)
			}

			log.Printf("Sent response for method %s", p.Req.Method)
		}
	}
}
