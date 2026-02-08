# MCP Server

A Model Context Protocol (MCP) server implementation in Go.

## Overview

Protocol version: `2025-03-26`

Supports:
- `initialize` - Server initialization
- `tools/list` - List available tools
- `tools/call` - Execute tool calls

## Tools

| Name | Description |
|------|-------------|
| `echo` | Echo back the provided text |

## Run

```bash
go run main.go
```

## TODO

- Handle multiple requests concurrently (currently single-threaded)
- Add more tools

## Notes

- Zero third-party dependencies (stdlib only)
- Communicates via stdin/stdout (stdio transport)
- Uses JSON-RPC 2.0 for message format
