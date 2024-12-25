// options.go
package go_jsonrpc

import (
	"log"
	"os"
)

// Options defines configuration options for JsRPC
type Options struct {
	CGI         bool   // Flag to control CGI header output
	Logger      Logger // Logger for logging events
	LogRequests bool   // Flag to log requests
}

// DefaultOptions provides default configuration for JsRPC
func DefaultOptions() *Options {
	return &Options{
		CGI:         true,
		Logger:      log.New(os.Stdout, "JsRPC: ", log.LstdFlags), // Default logger
		LogRequests: false,
	}
}
