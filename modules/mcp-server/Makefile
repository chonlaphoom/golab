
BINARY := mcp-server
TMPDIR := ./tmp

.PHONY: build build-race clean test-init test-discovery test-call test-all test-race test-race-only test-race-verbose test-race-report run

build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(TMPDIR)
	@go build -o $(BINARY) . && echo "Build successful."

build-race:
	@echo "Building $(BINARY) with race detector..."
	@mkdir -p $(TMPDIR)
	@go build -race -o $(BINARY) . && echo "Build with race detector successful."

clean:
	@rm -f $(BINARY)
	@rm -rf $(TMPDIR)

test-init: build
	@echo "=> initialize test"
	@mkdir -p $(TMPDIR)
	@printf '{"jsonrpc":"2.0","method":"initialize","params":[{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"mcp-client","version":"0.0.1"}}],"id":1}\n' | ./$(BINARY) > $(TMPDIR)/mcp_init_output.json || true
	@grep -q '"name":"mcp-server"' $(TMPDIR)/mcp_init_output.json && echo "initialize: OK" || (echo "initialize: FAIL"; cat $(TMPDIR)/mcp_init_output.json; exit 1)

test-discovery: build
	@echo "=> tools/list test"
	@mkdir -p $(TMPDIR)
	@printf '{"jsonrpc":"2.0","method":"tools/list","params":[],"id":2}\n' | ./$(BINARY) > $(TMPDIR)/mcp_list_output.json || true
	@grep -q '"echo"' $(TMPDIR)/mcp_list_output.json && echo "tools/list: OK" || (echo "tools/list: FAIL"; cat $(TMPDIR)/mcp_list_output.json; exit 1)

test-call: build
	@echo "=> tools/call (echo) test"
	@mkdir -p $(TMPDIR)
	@printf '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"echo","arguments":{"text":"hello world"}},"id":3}\n' | ./$(BINARY) > $(TMPDIR)/mcp_call_output.json || true
	@grep -q 'Echo: hello world' $(TMPDIR)/mcp_call_output.json && echo "tools/call: OK" || (echo "tools/call: FAIL"; cat $(TMPDIR)/mcp_call_output.json; exit 1)

test-jsonrpc: build
	@echo "=> JSON-RPC compliance test"
	@go test ./... -v | grep -q 'FAIL' && (echo "JSON-RPC compliance: FAIL"; exit 1) || echo "JSON-RPC compliance: OK"

test-all: test-init test-discovery test-call test-jsonrpc
	@echo "All tests passed."

test-race: build-race test-init test-discovery test-call test-jsonrpc
	@echo "All tests passed with race detector."

test-race-only:
	@echo "=> Running race condition tests only"
	@mkdir -p $(TMPDIR)
	@go test -race -run Race ./... -timeout 30s && echo "Race tests: PASS" || (echo "Race tests: FAIL"; exit 1)

test-race-verbose:
	@echo "=> Running race tests with verbose output"
	@mkdir -p $(TMPDIR)
	@go test -race -v -run Race ./... -timeout 30s

test-race-report:
	@echo "=> Running race tests with detailed report"
	@mkdir -p $(TMPDIR)
	@go test -race -v ./... -timeout 60s 2>&1 | tee $(TMPDIR)/race_report.txt
	@echo ""
	@echo "=== Race Detection Summary ==="
	@grep -q "WARNING: DATA RACE" $(TMPDIR)/race_report.txt && (echo "❌ RACE CONDITIONS DETECTED! Check $(TMPDIR)/race_report.txt for details"; exit 1) || echo "✅ No race conditions detected"

test-all-clean: clean test-all clean
	@echo "All tests passed and cleaned up." 

run:
	@go run .
