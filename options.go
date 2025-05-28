// options.go
package go_jsonrpc

import (
	"io"
	"log"
	"os"
)

// Options defines configuration options for JsRPC
type Options struct {
	CGI                bool   // Flag to control CGI header output
	Logger             Logger // Logger for logging events
	LogRequests        bool   // Flag to log requests
	HandlerInterceptor func(reader io.Reader, writer io.Writer) (finished bool, err error)
}

// DefaultOptions provides default configuration for JsRPC
func DefaultOptions() *Options {
	return &Options{
		CGI:         true,
		Logger:      log.New(os.Stdout, "JsRPC: ", log.LstdFlags), // Default logger
		LogRequests: false,
	}
}
