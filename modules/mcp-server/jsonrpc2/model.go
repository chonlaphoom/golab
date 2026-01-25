package jsonrpc2

type Message struct {
	Request         *Request
	SuccessResponse *SuccessResponse
	ErrorResponse   *ErrorResponse
}

type Request struct {
	JSONRPC_VER string `json:"jsonrpc"`
	Method      string `json:"method"`
	Params      []any  `json:"params"`
	ID          int    `json:"id"`
}

type SuccessResponse struct {
	JSONRPC_VER string      `json:"jsonrpc"`
	Result      interface{} `json:"result"`
	ID          int         `json:"id"`
}

type ErrorResponse struct {
	JSONRPC_VER string `json:"jsonrpc"`
	Error       Error  `json:"error"`
	ID          int    `json:"id"`
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    string    `json:"data,omitempty"`
}
