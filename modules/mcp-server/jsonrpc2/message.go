package jsonrpc2

import (
	"encoding/json"
	"errors"
)

func NewMessage() *Message {
	return &Message{
		Request: &Request{},
	}
}

func (e RPCError) Error() string { return e.Err.Error() }

// ParseRequest parses raw JSON request data into the message's Request and
// performs basic validation. It accepts both positional (array) and named
// (object) params; object params are wrapped into a single-element params
// array to preserve the original behavior.
func (msg *Message) ParseRequest(data []byte) error {
	var raw struct {
		JSONRPC json.RawMessage `json:"jsonrpc"`
		Method  string          `json:"method"`
		Params  json.RawMessage `json:"params"`
		ID      json.RawMessage `json:"id"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return RPCError{Code: Parse, Err: err}
	}

	// jsonrpc field must be the literal string "2.0"
	var ver string
	if len(raw.JSONRPC) > 0 {
		err := json.Unmarshal(raw.JSONRPC, &ver)
		if err != nil {
			return RPCError{Code: InvalidRequest, Err: errors.New("invalid jsonrpc field")}
		}
		if ver != "2.0" {
			return RPCError{Code: InvalidRequest, Err: errors.New("jsonrpc must be \"2.0\"")}
		}
	}

	// populate request
	if msg.Request == nil {
		msg.Request = &Request{}
	}
	msg.Request.JSONRPC = ver
	msg.Request.Method = raw.Method

	// parse params: array or object
	if len(raw.Params) > 0 {
		// detect leading character
		b := raw.Params
		// skip leading spaces
		i := 0
		for i < len(b) && (b[i] == ' ' || b[i] == '\n' || b[i] == '\t' || b[i] == '\r') {
			i++
		}
		if i < len(b) && b[i] == '[' {
			var arr []any
			if err := json.Unmarshal(raw.Params, &arr); err != nil {
				return RPCError{Code: InvalidParams, Err: err}
			}
			msg.Request.Params = arr
		} else {
			var obj map[string]any
			if err := json.Unmarshal(raw.Params, &obj); err != nil {
				return RPCError{Code: InvalidParams, Err: err}
			}
			msg.Request.Params = []any{obj}
		}
	} else {
		msg.Request.Params = []any{}
	}

	// parse id if present and numeric; otherwise leave as 0
	if len(raw.ID) > 0 {
		// try number
		var num float64
		if err := json.Unmarshal(raw.ID, &num); err == nil {
			msg.Request.ID = int(num)
		} else {
			// non-numeric IDs not supported by this implementation
			return RPCError{Code: InvalidRequest, Err: errors.New("unsupported id type; only numeric ids are supported")}
		}
	}

	// validate basic fields
	if msg.Request.JSONRPC != "2.0" {
		return RPCError{Code: InvalidRequest, Err: errors.New("jsonrpc must be \"2.0\"")}
	}
	if msg.Request.Method == "" {
		return RPCError{Code: InvalidRequest, Err: errors.New("method is required")}
	}

	return nil
}

func (msg *Message) WriteErrorResponse(id int, code ErrorCode, _message string) {
	var error_message string
	if _message != "" {
		error_message = _message
	} else {
		switch code {
		case Parse:
			error_message = `Invalid JSON was received by the server.
An error occurred on the server while parsing the JSON text.`
		case InvalidRequest:
			error_message = "The JSON sent is not a valid Request object."
		case MethodNotFound:
			error_message = "The method does not exist / is not available."
		case InvalidParams:
			error_message = "Invalid method parameter(s)."
		case InternalError:
			error_message = "Internal JSON-RPC error."
		default:
			error_message = "Unknown error."
		}
	}
	msg.ErrorResponse = &ErrorResponse{
		JSONRPC: "2.0",
		Error: Error{
			Code:    code,
			Message: error_message,
		},
		ID: msg.Request.ID,
	}
}

func (msg *Message) WriteSuccessResponse(result any) {
	msg.SuccessResponse = &SuccessResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      msg.Request.ID,
	}
}
