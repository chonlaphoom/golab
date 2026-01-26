package jsonrpc2

import (
	"testing"
)

func TestRegisterAndDispatch(t *testing.T) {
	// reset registry
	mu.Lock()
	registry = map[string]Handler{}
	mu.Unlock()

	RegisterFunc("test.sum", func(ctx *Context, params []any) (any, *ErrorResponse) {
		if len(params) < 2 {
			return nil, &ErrorResponse{JSONRPC: "2.0", Error: Error{Code: InvalidParams, Message: "need two params"}}
		}
		a, _ := params[0].(float64)
		b, _ := params[1].(float64)
		return a + b, nil
	})

	msg := NewMessage()
	msg.Request = &Request{JSONRPC: "2.0", Method: "test.sum", Params: []any{1.0, 2.0}, ID: 1}
	ctx := NewContext(msg)
	DispatchMessage(msg, ctx)

	if msg.SuccessResponse == nil {
		t.Fatalf("expected success response, got error: %+v", msg.ErrorResponse)
	}
	if msg.SuccessResponse.Result != 3.0 {
		t.Fatalf("expected 3.0, got %+v", msg.SuccessResponse.Result)
	}
}
