## Tiny Experiments Repository

Go concept experiments repository for learning and sharing small code snippets and projects.

"main.go" is for running experiments.

## Packages

### throttler/
Task throttling system with worker pool pattern. Features:
- Configurable worker count and buffer size
- Context-based cancellation
- Concurrent task execution with sync.WaitGroup
- Result channel for task completion tracking

### greeting/
Simple demonstration package showing basic Go package structure and imports.

### rune/
Demonstrates Unicode handling in Go:
- Difference between string byte length and rune count
- Converting strings to rune slices for proper character indexing
- Handling multi-byte UTF-8 characters (e.g., Chinese characters)
- Reliable character access using rune indexing

### append_/
Explores Go slice internals:
- Slice capacity and growth behavior
- Memory allocation when exceeding capacity
- Slice copying vs. referencing
- Using unsafe package to inspect memory addresses
