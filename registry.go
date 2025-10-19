// registry.go
package go_jsonrpc

import "os"

// command represents a registered command with its handler and specific middlewares.
type command struct {
	handler     HandlerFunc
	middlewares []MiddlewareFunc
}

// JsRPC is the main structure of the JSON-RPC server, handling registered commands and global middlewares.
type JsRPC struct {
	handlers    map[string]command
	middlewares []MiddlewareFunc // Global middlewares
	cgi         bool             // Flag to write CGI headers
	logger      Logger           // Logger for logging critical events
	options     *Options
	socketPerms os.FileMode
}

// HandlerFunc is the type definition for the function signature of a command handler.
type HandlerFunc func(ctx *Context) error

// MiddlewareFunc is the type definition for the function signature of a middleware.
type MiddlewareFunc func(ctx *Context) error

// New creates a new instance of JsRPC with the given options.
func New(options *Options) *JsRPC {
	if options == nil {
		options = DefaultOptions()
	}
	if options.Logger == nil {
		options.Logger = DefaultOptions().Logger
	}
	return &JsRPC{
		handlers: make(map[string]command),
		cgi:      options.CGI,
		logger:   options.Logger,
		options:  options,
	}
}

// RegisterCommand registers a command with a handler and optional middlewares.
func (r *JsRPC) RegisterCommand(commandName string, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	r.handlers[commandName] = command{
		handler:     handler,
		middlewares: middlewares,
	}
}

// UseGlobalMiddleware adds a global middleware that applies to all commands.
// Global middlewares are executed before command-specific middlewares.
// If a global middleware returns an error, the execution of the command is stopped.
func (r *JsRPC) UseGlobalMiddleware(middleware MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middleware)
}
