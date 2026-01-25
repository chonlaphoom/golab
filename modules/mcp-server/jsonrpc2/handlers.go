package jsonrpc2

import "sync"

type Handler interface {
	Handle(ctx *Context, params []any) (any, *ErrorResponse)
}

type HandlerFunc func(ctx *Context, params []any) (any, *ErrorResponse)

func (f HandlerFunc) Handle(ctx *Context, params []any) (any, *ErrorResponse) {
	return f(ctx, params)
}

var (
	mu       sync.RWMutex
	registry = map[string]Handler{}
)

func register(method string, h Handler) {
	mu.Lock()
	defer mu.Unlock()
	registry[method] = h
}

func RegisterFunc(method string, fn HandlerFunc) {
	register(method, fn)
}

func Unregister(method string) {
	mu.Lock()
	defer mu.Unlock()
	delete(registry, method)
}

func getHandler(method string) Handler {
	mu.RLock()
	defer mu.RUnlock()
	return registry[method]
}

// DispatchMessage looks up the handler for msg.Request.Method, executes it, and
// fills msg.SuccessResponse or msg.ErrorResponse accordingly.
func DispatchMessage(msg *Message, ctx *Context) {
	if msg == nil || msg.Request == nil {
		msg.NewErrorResponse(0, InvalidRequest, "")
		return
	}

	h := getHandler(msg.Request.Method)
	if h == nil {
		msg.NewErrorResponse(msg.Request.ID, MethodNotFound, "")
		return
	}

	result, errResp := h.Handle(ctx, msg.Request.Params)
	if errResp != nil {
		if errResp.JSONRPC_VER == "" {
			errResp.JSONRPC_VER = "2.0"
		}
		errResp.ID = msg.Request.ID
		msg.ErrorResponse = errResp
		return
	}

	msg.NewSuccessResponseWithResult(result)
}
