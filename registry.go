// registry.go
package go_jsonrpc

type command struct {
	handler     HandlerFunc
	middlewares []MiddlewareFunc
}

type JsRPC struct {
	handlers    map[string]command
	middlewares []MiddlewareFunc // Global middlewares
	cgi         bool             // Flag to write CGI headers
}

type HandlerFunc func(ctx *Context)
type MiddlewareFunc func(ctx *Context) error

// New creates a new instance of JsRPC with a flag to control CGI header output
func New(cgi bool) *JsRPC {
	return &JsRPC{
		handlers: make(map[string]command),
		cgi:      cgi,
	}
}

// RegisterCommand registers a command with optional middlewares
func (r *JsRPC) RegisterCommand(commandName string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.handlers[commandName] = command{
		handler:     handler,
		middlewares: middlewares,
	}
}

// UseGlobalMiddleware adds a middleware that applies to all commands
func (r *JsRPC) UseGlobalMiddleware(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}
