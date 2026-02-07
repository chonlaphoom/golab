package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const protocolVersion = "2026-01-26"

// func handleInitialize(request mcp.InitializeRequest) mcp.InitializeResponse {
// 	return mcp.InitializeResponse{
// 		ProtocolVersion: protocolVersion,
// 		Capabilities: mcp.ServerCapabilities{
// 			Tools: map[string]any{
// 				"listChanged": false, // tools does not change during the session
// 			},
// 		},
// 		ServerInfo: mcp.ServerInfo{
// 			Name:    "demo-mcp-server",
// 			Version: "0.0.1",
// 		},
// 	}
// }
//
// func handleListTools() mcp.ListToolsResponse {
// 	return mcp.ListToolsResponse{
// 		Tools: []mcp.Tool{
// 			{
// 				Name:        "echo",
// 				Description: "Echo back the provided text",
// 				InputSchema: map[string]any{
// 					"type": "object",
// 					"properties": map[string]any{
// 						"text": map[string]any{
// 							"type":        "string",
// 							"description": "Text to echo back",
// 						},
// 					},
// 					"required": []string{"text"},
// 				},
// 			},
// 		},
// 	}
// }
//
// func handleCallTool(call mcp.CallToolRequest) mcp.ToolResult {
// 	switch call.Name {
// 	case "echo":
// 		if text, ok := call.Arguments["text"].(string); ok {
// 			return mcp.ToolResult{
// 				Content: []any{
// 					map[string]any{
// 						"type": "text",
// 						"text": fmt.Sprintf("Echo: %s", text),
// 					},
// 				},
// 			}
// 		}
// 		return mcp.ToolResult{
// 			Content: []any{
// 				map[string]any{
// 					"type": "text",
// 					"text": "Error: 'text' argument is required and must be a string",
// 				},
// 			},
// 			IsError: true,
// 		}
// 	default:
// 		return mcp.ToolResult{
// 			Content: []any{
// 				map[string]any{
// 					"type": "text",
// 					"text": fmt.Sprintf("Unknown tool: %s", call.Name),
// 				},
// 			},
// 			IsError: true,
// 		}
// 	}
// }

func main() {
	log.Printf("MCP Server %s is running...", protocolVersion)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	decoder := json.NewDecoder(os.Stdin)

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
		}
	}
}
