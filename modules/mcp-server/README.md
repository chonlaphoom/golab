# MCP Server

A Model Context Protocol (MCP) server implementation in Go with concurrent request handling.

## Overview

Protocol version: `2025-03-26`

### Features

- **Concurrent Processing**: 5 worker goroutines handle requests in parallel
- **Thread-Safe Output**: Mutex-protected writer prevents output corruption
- **Graceful Shutdown**: Workers drain pending messages on context cancellation
- **Zero Dependencies**: Pure Go stdlib implementation

### Supported Methods

- `initialize` - Server initialization
- `tools/list` - List available tools
- `tools/call` - Execute tool calls

## Tools

| Name   | Description                 |
| ------ | --------------------------- |
| `echo` | Echo back the provided text |

## Building

```bash
# Standard build
make build

# Build with race detector
make build-race
```

## Running

```bash
# Run directly
go run .

# Or use make
make run
```

## Testing

### Run All Tests

```bash
make test-all
```

### Run Race Condition Tests

The server includes comprehensive race condition tests to verify concurrent safety:

```bash
# Run only race tests
make test-race-only

# Run with verbose output
make test-race-verbose

# Generate detailed race report
make test-race-report
```

**Test Coverage:**

- ✅ ThreadSafeWriter mutex synchronization (4 tests)
- ✅ Worker pool concurrent processing (6 tests)
- ✅ Message reader/distributor (7 tests)
- ✅ 17 total race condition tests
- ✅ Medium load testing (100-500 messages per test)

See [RACE_TESTS.md](./RACE_TESTS.md) for detailed documentation.

### Integration Tests

```bash
# Test server initialization
make test-init

# Test tool discovery
make test-discovery

# Test tool execution
make test-call

# Run all tests with cleanup
make test-all-clean
```

## Architecture

### Concurrency Model

```
stdin → readAndPushMsgs → msgChan → [Worker 1]
                                   → [Worker 2]
                                   → [Worker 3] → ThreadSafeWriter → stdout
                                   → [Worker 4]
                                   → [Worker 5]
```

- **Reader Goroutine**: Reads JSON-RPC messages from stdin, pushes to channel
- **Worker Pool**: 5 workers consume from message channel
- **Thread-Safe Writer**: Serializes concurrent writes to stdout with mutex
- **Graceful Shutdown**: Context cancellation triggers message draining

### Key Components

- `main.go` - Server initialization and coordination
- `worker.go` - Worker pool and message processing
- `read_and_distribute.go` - Message reader
- `jsonrpc2/` - JSON-RPC 2.0 parser and response builder
- `mcp/` - MCP protocol handlers

## Development

### Adding New Tools

1. Define tool schema in `mcp/handler.go` → `HandleListTools`
2. Add case in `mcp/handler.go` → `HandleCallTool`
3. Implement tool logic
4. Add tests

### Clean Build

```bash
make clean
make build
```

## Notes

- **Zero third-party dependencies** (stdlib only)
- **Communicates via stdin/stdout** (stdio transport)
- **Uses JSON-RPC 2.0** for message format
- **Thread-safe** concurrent request handling

## TODO

- [ ] Add timeout for worker tasks
- [ ] Add more tools
- [ ] Add prometheus metrics
- [ ] Add request tracing
