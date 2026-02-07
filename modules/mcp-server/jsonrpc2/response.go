package jsonrpc2

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Data    any       `json:"data,omitempty"`
}

type Response[T any] struct {
	JSONRPC string `json:"jsonrpc"`
	Result  T      `json:"result,omitempty"`
	Error   *Error `json:"error,omitempty"`
	ID      int    `json:"id"`
}

func NewSuccess[T any](id int, result T) Response[T] {
	return Response[T]{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	}
}

func NewError(id int, code ErrorCode, message string, data any) Response[any] {
	return Response[any]{
		JSONRPC: "2.0",
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	}
}
