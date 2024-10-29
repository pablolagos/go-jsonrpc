package go_jsonrpc

type Command struct {
	Handler     HandlerFunc
	Middlewares []MiddlewareFunc
}

type JsRPC struct {
	handlers    map[string]Command
	middlewares []MiddlewareFunc // Global middlewares
}

type HandlerFunc func(ctx *Context) // No error return value

type MiddlewareFunc func(ctx *Context) error

func NewRegistry() *JsRPC {
	return &JsRPC{
		handlers: make(map[string]Command),
	}
}

// RegisterCommand allows registering a command with optional middlewares
func (r *JsRPC) RegisterCommand(command string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.handlers[command] = Command{
		Handler:     handler,
		Middlewares: middlewares,
	}
}

// UseGlobalMiddleware adds a middleware that applies to all commands
func (r *JsRPC) UseGlobalMiddleware(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}
