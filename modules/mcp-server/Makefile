
BINARY := mcp-server
TMPDIR := ./tmp

.PHONY: build clean test-init test-discovery test-call test-all run

build:
	@echo "Building $(BINARY)..."
	@mkdir -p $(TMPDIR)
	@go build -o $(BINARY) . && echo "Build successful."

clean:
	@rm -f $(BINARY)
	@rm -rf $(TMPDIR)

test-init: build
	@echo "=> initialize test"
	@mkdir -p $(TMPDIR)
	@printf '{"jsonrpc":"2.0","method":"initialize","params":[{"protocolVersion":"2026-01-26","capabilities":{},"clientInfo":{"name":"mcp-client","version":"0.0.1"}}],"id":1}\n' | ./$(BINARY) > $(TMPDIR)/mcp_init_output.json || true
	@grep -q '"name":"demo-mcp-server"' $(TMPDIR)/mcp_init_output.json && echo "initialize: OK" || (echo "initialize: FAIL"; cat $(TMPDIR)/mcp_init_output.json; exit 1)

test-discovery: build
	@echo "=> tools/list test"
	@mkdir -p $(TMPDIR)
	@printf '{"jsonrpc":"2.0","method":"tools/list","params":[],"id":2}\n' | ./$(BINARY) > $(TMPDIR)/mcp_list_output.json || true
	@grep -q '"echo"' $(TMPDIR)/mcp_list_output.json && echo "tools/list: OK" || (echo "tools/list: FAIL"; cat $(TMPDIR)/mcp_list_output.json; exit 1)

test-call: build
	@echo "=> tools/call (echo) test"
	@mkdir -p $(TMPDIR)
	@printf '{"jsonrpc":"2.0","method":"tools/call","params":[{"name":"echo","arguments":{"text":"hello world"}}],"id":3}\n' | ./$(BINARY) > $(TMPDIR)/mcp_call_output.json || true
	@grep -q 'Echo: hello world' $(TMPDIR)/mcp_call_output.json && echo "tools/call: OK" || (echo "tools/call: FAIL"; cat $(TMPDIR)/mcp_call_output.json; exit 1)

test-all: test-init test-discovery test-call
	@echo "All tests passed."

test-all-clean: clean test-all clean
	@echo "All tests passed and cleaned up." 

run:
	@go run .
