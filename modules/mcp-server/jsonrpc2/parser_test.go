package jsonrpc2

import (
	"testing"
)

func TestParser_ParseRequest_InvalidJSON(t *testing.T) {
	parser := NewParser()
	invalidJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": [1, 2, 3], "id": 1`)
	err := parser.ParseRequest(invalidJSON)
	if err == nil {
		t.Errorf("Expected error for invalid JSON, got nil")
	}
}

func TestParser_ParseRequest_MissingMethod(t *testing.T) {
	parser := NewParser()
	missingMethodJSON := []byte(`{"jsonrpc": "2.0", "params": [1, 2, 3], "id": 1}`)
	err := parser.ParseRequest(missingMethodJSON)
	if err == nil {
		t.Errorf("Expected error for missing method, got nil")
	}
}

func TestParser_ParseRequest_ValidRequest(t *testing.T) {
	parser := NewParser()
	validJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": [1, 2, 3], "id": 1}`)
	err := parser.ParseRequest(validJSON)
	if err != nil {
		t.Errorf("Unexpected error for valid request: %v", err)
	}
	if parser.Req.Method != "testMethod" {
		t.Errorf("Expected method 'testMethod', got '%s'", parser.Req.Method)
	}
	if parser.Req.ID != 1 {
		t.Errorf("Expected ID 1, got %d", parser.Req.ID)
	}
}

func TestParser_ParseRequest_Notification(t *testing.T) {
	parser := NewParser()
	notificationJSON := []byte(`{"jsonrpc": "2.0", "method": "notifyMethod", "params": [1, 2, 3]}`)
	err := parser.ParseRequest(notificationJSON)
	if err != nil {
		t.Errorf("Unexpected error for notification: %v", err)
	}
	if !parser.Req.IsNotification {
		t.Errorf("Expected IsNotification to be true, got false")
	}
}

func TestParser_ParseRequest_InvalidID(t *testing.T) {
	parser := NewParser()
	invalidIDJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": [1, 2, 3], "id": "stringID"}`)
	err := parser.ParseRequest(invalidIDJSON)
	if err == nil {
		t.Errorf("Expected error for invalid ID type, got nil")
	}
}

func TestParser_ParseRequest_FloatID(t *testing.T) {
	parser := NewParser()
	floatIDJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": [1, 2, 3], "id": 1.5}`)
	err := parser.ParseRequest(floatIDJSON)
	if err != nil {
		t.Errorf("Unexpected error for float ID: %v", err)
	}
	if parser.Req.ID != 1 {
		t.Errorf("Expected ID 1, got %d", parser.Req.ID)
	}
}

func TestParser_ParseRequest_NullID(t *testing.T) {
	parser := NewParser()
	nullIDJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": [1, 2, 3], "id": null}`)
	err := parser.ParseRequest(nullIDJSON)
	if err != nil {
		t.Errorf("Unexpected error for null ID: %v", err)
	}
	if parser.Req.ID != 0 {
		t.Errorf("Expected ID 0 for null, got %d", parser.Req.ID)
	}
}

func TestParser_ParseRequest_ObjectParams(t *testing.T) {
	parser := NewParser()
	objectParamsJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": {"key1": "value1", "key2": 2}, "id": 1}`)
	err := parser.ParseRequest(objectParamsJSON)
	if err != nil {
		t.Errorf("Unexpected error for object params: %v", err)
	}
	if param, ok := parser.Req.Params.data.(map[string]any); !ok {
		t.Errorf("Expected Params to be map[string]any, got %T", parser.Req.Params.data)
	} else if param["key1"] != "value1" || param["key2"].(float64) != 2 {
		t.Errorf("Unexpected Params content: %v", param)
	}
}

func TestParser_ParseRequest_ArrayParams(t *testing.T) {
	parser := NewParser()
	arrayParamsJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "params": [1, 2, 3], "id": 1}`)
	err := parser.ParseRequest(arrayParamsJSON)
	if err != nil {
		t.Errorf("Unexpected error for array params: %v", err)
	}
	if param, ok := parser.Req.Params.data.([]any); !ok {
		t.Errorf("Expected Params to be []any, got %T", parser.Req.Params.data)
	} else {
		if len(param) != 3 {
			t.Errorf("Expected Params length 3, got %d", len(param))
		}
		if param[0].(float64) != 1.0 {
			t.Errorf("Expected Params[0] to be 1, got %v", param[0])
		}
		if param[1].(float64) != 2.0 {
			t.Errorf("Expected Params[1] to be 2, got %v", param[1])
		}
		if param[2].(float64) != 3.0 {
			t.Errorf("Expected Params[2] to be 3, got %v", param[2])
		}
	}
}

func TestParser_parseRequest_EmptyParams(t *testing.T) {
	parser := NewParser()
	noParamsJSON := []byte(`{"jsonrpc": "2.0", "method": "testMethod", "id": 1}`)
	err := parser.ParseRequest(noParamsJSON)
	if err != nil {
		t.Errorf("Unexpected error for no params: %v", err)
	}

	if parser.Req.Params.data != nil {
		t.Errorf("Expected Params to be nil, got %v", parser.Req.Params.data)
	}
}
