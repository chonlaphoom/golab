package jsonrpc2

type RPCError struct {
	Code ErrorCode
	Err  error
}

func (e RPCError) Error() string {
	switch e.Code {
	case Parse:
		return "Invalid JSON was received by the server. An error occurred on the server while parsing the JSON text."
	case InvalidRequest:
		return "The JSON sent is not a valid Request object."
	case MethodNotFound:
		return "The method does not exist / is not available."
	case InvalidParams:
		return "Invalid method parameter(s)."
	case InternalError:
		return "Internal JSON-RPC error."
	default:
		return "Unknown error."
	}
}

type ErrorCode int

const Parse ErrorCode = -32700
const InvalidRequest ErrorCode = -32600
const MethodNotFound ErrorCode = -32601
const InvalidParams ErrorCode = -32602
const InternalError ErrorCode = -32603
