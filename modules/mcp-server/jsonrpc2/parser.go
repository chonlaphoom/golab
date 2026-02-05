package jsonrpc2

import (
	"encoding/json"
	"fmt"
)

type Parser struct {
	Req Request
}

type Request struct {
	JSONRPC        string
	Method         string
	Params         any
	ID             int
	IsNotification bool
}

func NewParser() *Parser {
	return &Parser{}
}

/* parse and validate a request */
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

	if p.Req.JSONRPC != "2.0" {
		return RPCError{Code: InvalidRequest, Err: fmt.Errorf("jsonrpc must be \"2.0\"")}
	}
	p.Req.JSONRPC = raw.Jsonrpc

	if p.Req.Method == "" {
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
			return RPCError{Code: InvalidRequest, Err: fmt.Errorf("string IDs are not supported")}
		case nil:
		default:
			p.Req.IsNotification = true
		}
	}

	// parse params

	return nil
}
