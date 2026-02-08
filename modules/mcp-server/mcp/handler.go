package mcp

import (
	"fmt"
	"mcp-server/jsonrpc2"
)

func HandleInitialize(ver string, id int) jsonrpc2.Response[InitializeResponse] {
	res := InitializeResponse{
		ProtocolVersion: ver,
		Capabilities: ServerCapabilities{
			Tools: map[string]any{
				"listChanged": false, // tools does not change during the session
			},
		},
		ServerInfo: ServerInfo{
			Name:    "mcp-server",
			Version: "0.0.1",
		},
	}
	return jsonrpc2.NewSuccess(id, res)
}

func HandleListTools(id int) jsonrpc2.Response[ListToolsResponse] {
	return jsonrpc2.NewSuccess(id, ListToolsResponse{
		Tools: []Tool{
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
	})
}

func HandleCallTool(params jsonrpc2.Params, id int) jsonrpc2.Response[ToolResult] {
	obj, isObj := params.GetAsObject()

	call := CallToolRequest{
		Name:      "",
		Arguments: map[string]any{},
	}

	if !isObj {
		return jsonrpc2.NewError[ToolResult](id, jsonrpc2.InvalidParams, "", ToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": "Error: params must be an object",
				},
			},
			IsError: true,
		})
	}

	if n, ok := obj["name"].(string); ok {
		call.Name = n
	}
	if args, ok := obj["arguments"].(map[string]any); ok {
		call.Arguments = args
	}

	switch call.Name {
	case "echo":
		if text, ok := call.Arguments["text"].(string); ok {
			return jsonrpc2.NewSuccess(id, ToolResult{
				Content: []any{
					map[string]any{
						"type": "text",
						"text": fmt.Sprintf("Echo: %s", text),
					},
				},
			})
		}
		return jsonrpc2.NewError[ToolResult](id, jsonrpc2.InvalidParams, "", ToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": "Error: 'text' argument is required and must be a string",
				},
			},
			IsError: true,
		})
	default:
		return jsonrpc2.NewError[ToolResult](id, jsonrpc2.InvalidParams, "", ToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": fmt.Sprintf("Unknown tool: %s", call.Name),
				},
			},
			IsError: true,
		})
	}
}
