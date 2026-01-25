package jsonrpc2

import (
	"encoding/json"
)

func NewMessage() *Message {
	return &Message{
		Request: &Request{},
	}
}

func (req *Message) ParseRequest(data []byte) error {
	err := json.Unmarshal(data, req.Request)
	if err != nil {
		return err
	}

	return nil
}

type ErrorCode int

const Parse ErrorCode = -32700
const InvalidRequest ErrorCode = -32600
const MethodNotFound ErrorCode = -32601
const InvalidParams ErrorCode = -32602
const InternalError ErrorCode = -32603

func (msg *Message) NewErrorResponse(id int, code ErrorCode, _message string) {
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
		JSONRPC_VER: "2.0",
		Error: Error{
			Code:    code,
			Message: error_message,
		},
		ID: msg.Request.ID,
	}
}

func (msg *Message) NewSuccessResponseWithResult(result any) {
	msg.SuccessResponse = &SuccessResponse{
		JSONRPC_VER: "2.0",
		Result:      result,
		ID:          msg.Request.ID,
	}
}

func (succ *SuccessResponse) UnmarshalJSON() ([]byte, error) {
	b, err := json.Marshal(succ)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
