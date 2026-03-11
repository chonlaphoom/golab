# GoLab

GoLab is a collection of experiments and modules

Each experiment or module is organized into packages for easy management and execution.

## Structure

### pkgs/
Small experiment packages for learning Go concepts:

- **throttler/** - Task throttler with worker pool pattern, context cancellation, and concurrent task execution
- **greeting/** - Simple greeting package demonstrating basic package structure
- **rune/** - Demonstrates Go rune (Unicode code point) handling, UTF-8 encoding, and string/rune slice indexing
- **append_/** - Explores Go slice append behavior, capacity management, and memory allocation

### modules/
Larger experiment modules with multiple packages:

- **crawler/** - Experimental web crawler with concurrent page fetching, URL normalization, and CSV export
- **mcp-server/** - Model Context Protocol (MCP) server with JSON-RPC 2.0, concurrent worker pool (5 workers), and thread-safe output handling
- **worker-pool/** - Worker pool implementation example
- **rpc/** - RPC experiments including gRPC and Twirp implementations
