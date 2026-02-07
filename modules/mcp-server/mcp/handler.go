package mcp

import "fmt"

func HandleInitialize(ver string) InitializeResponse {
	return InitializeResponse{
		ProtocolVersion: ver,
		Capabilities: ServerCapabilities{
			Tools: map[string]any{
				"listChanged": false, // tools does not change during the session
			},
		},
		ServerInfo: ServerInfo{
			Name:    "demo-server",
			Version: "0.0.1",
		},
	}
}

func HandleListTools() ListToolsResponse {
	return ListToolsResponse{
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
	}
}

func HandleCallTool(call CallToolRequest) ToolResult {
	switch call.Name {
	case "echo":
		if text, ok := call.Arguments["text"].(string); ok {
			return ToolResult{
				Content: []any{
					map[string]any{
						"type": "text",
						"text": fmt.Sprintf("Echo: %s", text),
					},
				},
			}
		}
		return ToolResult{
			Content: []any{
				map[string]any{
					"type": "text",
					"text": "Error: 'text' argument is required and must be a string",
				},
			},
			IsError: true,
		}
	default:
		return ToolResult{
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
