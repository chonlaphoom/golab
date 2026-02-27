package jsonrpc2

type Error struct {
	Data    any       `json:"data,omitempty"`
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

type Response[T any] struct {
	Result  T      `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

func NewSuccess[T any](id int, result T) Response[T] {
	return Response[T]{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
}

func NewError[T any](id int, code ErrorCode, message string, data any) Response[T] {
	return Response[T]{
		JSONRPC: "2.0",
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}
}
