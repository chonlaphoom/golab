package jsonrpc2

import (
	"encoding/json"
	"fmt"
)

type Parser struct {
	Req Request
	Res Response[any]
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseRequest(data []byte) error {
	type rawRequest struct {
		Jsonrpc string `json:"jsonrpc"`
		Method  string `json:"method"`
		// array or object
		Params json.RawMessage `json:"params"`
		// string, number, or null
		Id json.RawMessage `json:"id"`
	}

	var raw rawRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return RPCError{Code: Parse, Err: err}
	}

	if raw.Jsonrpc != "2.0" {
		return RPCError{Code: InvalidRequest, Err: fmt.Errorf("jsonrpc must be \"2.0\"")}
	}
	p.Req.JSONRPC = raw.Jsonrpc

	if raw.Method == "" {
		return RPCError{Code: InvalidRequest, Err: fmt.Errorf("method is required")}
	}
	p.Req.Method = raw.Method

	var temp any
	if err := json.Unmarshal(raw.Id, &temp); err != nil {
		// malformed id, assume notification
		p.Req.IsNotification = true
	} else {
		switch v := temp.(type) {
		case float64:
			p.Req.ID = int(v)
		case string:
			return RPCError{Code: InvalidRequest, Err: fmt.Errorf("string IDs are not supported, for simplicity")}
		case nil:
		default:
			p.Req.IsNotification = true
		}
	}

	if len(raw.Params) > 0 {
		// skip white space
		index := 0
		for index < len(raw.Params) && (raw.Params[index] == ' ' || raw.Params[index] == '\n' || raw.Params[index] == '\t' || raw.Params[index] == '\r') {
			index++
		}

		if index < len(raw.Params) {
			switch raw.Params[index] {
			case '[':
				var arrParams []any
				err := json.Unmarshal(raw.Params, &arrParams)
				if err != nil {
					return RPCError{Code: InvalidParams, Err: err}
				}

				// TODO: check if all elements are same type
				p.Req.Params.data = arrParams
			case '{':
				var objParams map[string]any
				if err := json.Unmarshal(raw.Params, &objParams); err != nil {
					return RPCError{Code: InvalidParams, Err: err}
				}
				p.Req.Params.data = objParams
			default:
				return RPCError{Code: InvalidParams, Err: fmt.Errorf("params must be array or object")}
			}
		}
	}

	return nil
}
