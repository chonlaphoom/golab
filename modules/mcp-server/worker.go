package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mcp-server/jsonrpc2"
	"mcp-server/mcp"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, safeWriter *ThreadSafeWriter, msgChan <-chan json.RawMessage, errChan <-chan error, work_id int) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println(work_id, ": Received shutdown signal, draining remaining messages...")
			// Drain phase: process remaining messages until channel closes
			for msg := range msgChan {
				processMessage(safeWriter, msg)
			}
			log.Println(work_id, ": Message channel drained, exiting.")
			return
		case err := <-errChan:
			if errors.Is(err, io.EOF) {
				log.Println(work_id, ": EOF received, exiting.")
				return
			}
			log.Printf("%d : Error reading message: %v", work_id, err)
			return
		case msg, ok := <-msgChan:
			if !ok {
				log.Println(work_id, ": Message channel closed, exiting.")
				return
			}
			processMessage(safeWriter, msg)
		}
	}
}

func processMessage(safeWriter *ThreadSafeWriter, msg json.RawMessage) {
	log.Printf("Message received on channel: %s", string(msg))

	log.Printf("Received message: %s", string(msg))
	p := jsonrpc2.NewParser()
	if err := p.ParseRequest(msg); err != nil {
		log.Printf("Error parsing request: %v, wait for next message", err)
		return
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
	case "notifications/initialized":
		log.Println("Client initialized notification received")
		return
	case "tools/list":
		res = mcp.HandleListTools(p.Req.ID)
	case "tools/call":
		res = mcp.HandleCallTool(p.Req.Params, p.Req.ID)
	default:
		err = errors.New("unknown method: " + p.Req.Method)
	}
	if err != nil {
		log.Printf("Error handling method %s: %v", p.Req.Method, err)
		return
	}

	safeWriter.mu.Lock()
	defer safeWriter.mu.Unlock()

	en := json.NewEncoder(safeWriter.writer)
	if err := en.Encode(res); err != nil {
		log.Printf("Error encoding response: %v", err)
		return
	}

	if err := safeWriter.writer.Flush(); err != nil {
		log.Printf("Error flushing writer: %v", err)
		return
	}

	log.Printf("Sent response for method %s", p.Req.Method)
}
