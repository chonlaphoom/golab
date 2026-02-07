package jsonrpc2

type Request struct {
	JSONRPC        string
	Method         string
	Params         Params
	ID             int
	IsNotification bool
}

type Params struct {
	data any
}

func (p *Params) IsArray() bool {
	_, ok := p.data.([]any)
	return ok
}

func (p *Params) IsObject() bool {
	_, ok := p.data.(map[string]any)
	return ok
}

func (p *Params) GetAsArray() ([]any, bool) {
	arr, ok := p.data.([]any)
	return arr, ok
}

func (p *Params) GetAsObject() (map[string]any, bool) {
	obj, ok := p.data.(map[string]any)
	return obj, ok
}
