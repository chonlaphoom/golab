package jsonrpc2

import "testing"

func TestResponse_NewSuccess(t *testing.T) {
	id := 1
	result := "success"
	resp := NewSuccess(id, result)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", resp.JSONRPC)
	}
	if resp.Result != result {
		t.Errorf("Expected Result '%s', got '%s'", result, resp.Result)
	}
	if resp.ID != id {
		t.Errorf("Expected ID %d, got %d", id, resp.ID)
	}
	if resp.Error != nil {
		t.Errorf("Expected Error to be nil, got %v", resp.Error)
	}
}

func TestResponse_NewError(t *testing.T) {
	id := 1
	code := InvalidRequest
	message := "Invalid request"
	data := "Additional error data"
	resp := NewError[string](id, code, message, data)

	if resp.JSONRPC != "2.0" {
		t.Errorf("Expected JSONRPC '2.0', got '%s'", resp.JSONRPC)
	}
	if resp.ID != id {
		t.Errorf("Expected ID %d, got %d", id, resp.ID)
	}
	if resp.Error == nil {
		t.Errorf("Expected Error to be non-nil, got nil")
	} else {
		if resp.Error.Code != code {
			t.Errorf("Expected Error Code %d, got %d", code, resp.Error.Code)
		}
		if resp.Error.Message != message {
			t.Errorf("Expected Error Message '%s', got '%s'", message, resp.Error.Message)
		}
		if resp.Error.Data != data {
			t.Errorf("Expected Error Data '%v', got '%v'", data, resp.Error.Data)
		}
	}
}
