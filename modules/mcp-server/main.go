package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"mcp-server/jsonrpc2"
	"mcp-server/mcp"
)

func handleInitialize(request mcp.InitializeRequest) mcp.InitializeResponse {
	return mcp.InitializeResponse{
		ProtocolVersion: "2024-11-05",
		Capabilities: mcp.ServerCapabilities{
			Tools: map[string]any{
				"listChanged": false,
			},
		},
		ServerInfo: mcp.ServerInfo{
			Name:    "demo-mcp-server",
			Version: "1.0.0",
		},
	}
}

func handleListTools() mcp.ListToolsResponse {
	return mcp.ListToolsResponse{
		Tools: []mcp.Tool{
			{
				Name:        "echo",
				Description: "Echo back the provided text",
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"text": map[string]any{
							"type":        "string",
							"description": "Text to echo back",
						},
					},
					"required": []string{"text"},
				},
			},
		},
	}
}

func handleCallTool(call mcp.CallToolRequest) mcp.ToolResult {
	switch call.Name {
	case "echo":
		if text, ok := call.Arguments["text"].(string); ok {
			return mcp.ToolResult{
				Content: []any{
					map[string]any{
						"type": "text",
						"text": fmt.Sprintf("Echo: %s", text),
					},
				},
			}
		}
		return mcp.ToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": "Error: 'text' argument is required and must be a string",
				},
			},
			IsError: true,
		}
	default:
		return mcp.ToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": fmt.Sprintf("Unknown tool: %s", call.Name),
				},
			},
			IsError: true,
		}
	}
}

func main() {
	jsonrpc2.RegisterFunc("initialize", func(ctx *jsonrpc2.Context, params []any) (any, *jsonrpc2.ErrorResponse) {
		var initReq mcp.InitializeRequest
		if len(params) > 0 {
			if m, ok := params[0].(map[string]any); ok {
				if b, err := json.Marshal(m); err == nil {
					_ = json.Unmarshal(b, &initReq)
				}
			}
		}
		return handleInitialize(initReq), nil
	})

	jsonrpc2.RegisterFunc("tools/list", func(ctx *jsonrpc2.Context, params []any) (any, *jsonrpc2.ErrorResponse) {
		return handleListTools(), nil
	})

	jsonrpc2.RegisterFunc("tools/call", func(ctx *jsonrpc2.Context, params []any) (any, *jsonrpc2.ErrorResponse) {
		var callReq mcp.CallToolRequest
		if len(params) > 0 {
			if m, ok := params[0].(map[string]any); ok {
				if b, err := json.Marshal(m); err == nil {
					_ = json.Unmarshal(b, &callReq)
				}
			}
		}
		return handleCallTool(callReq), nil
	})

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	decoder := json.NewDecoder(os.Stdin)

	for {
		select {
		case <-sigChan:
			slog.Info("Received shutdown signal, gracefully closing...")
			return
		default:
		}

		// Decode the next JSON value as raw bytes and let jsonrpc2 parse it
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			slog.Error("Error parsing JSON-RPC request", "err", err)
			continue
		}

		msg := jsonrpc2.NewMessage()
		if err := msg.ParseRequest(raw); err != nil {
			// parsing error -> respond with Parse error
			msg.NewErrorResponse(0, jsonrpc2.Parse, err.Error())
			out, _ := json.Marshal(msg.ErrorResponse)
			os.Stdout.Write(out)
			os.Stdout.WriteString("\n")
			os.Stdout.Sync()
			continue
		}

		// Dispatch using jsonrpc2
		ctx := jsonrpc2.NewContext(msg)
		jsonrpc2.DispatchMessage(msg, ctx)

		var out any
		if msg.ErrorResponse != nil {
			out = msg.ErrorResponse
		} else if msg.SuccessResponse != nil {
			out = msg.SuccessResponse
		} else {
			msg.NewErrorResponse(msg.Request.ID, jsonrpc2.InternalError, "no response")
			out = msg.ErrorResponse
		}

		responseBytes, err := json.Marshal(out)
		if err != nil {
			slog.Error("Error marshaling response", "err", err)
			continue
		}

		os.Stdout.Write(responseBytes)
		os.Stdout.WriteString("\n")
		os.Stdout.Sync()
	}
}
