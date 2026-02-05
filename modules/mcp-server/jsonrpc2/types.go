package jsonrpc2

type Message struct {
	Request         *Request
	SuccessResponse *SuccessResponse
	ErrorResponse   *ErrorResponse
}

type SuccessResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	ID      int         `json:"id"`
}

type ErrorResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Error   Error  `json:"error"`
	ID      int    `json:"id"`
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    string    `json:"data,omitempty"`
}

type RPCError struct {
	Code ErrorCode
	Err  error
}

type ErrorCode int

const Parse ErrorCode = -32700
const InvalidRequest ErrorCode = -32600
const MethodNotFound ErrorCode = -32601
const InvalidParams ErrorCode = -32602
const InternalError ErrorCode = -32603
