package jsonrpc2

type Context struct {
	Message *Message
	Request *Request

	ConnInfo map[string]any

	Values map[string]any
}

func NewContext(msg *Message) *Context {
	ctx := &Context{
		Message:  msg,
		Values:   map[string]any{},
		ConnInfo: map[string]any{},
	}
	if msg != nil {
		ctx.Request = msg.Request
	}
	return ctx
}
