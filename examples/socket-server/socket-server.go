package main

import (
	"errors"
	"log"
	"os"

	"github.com/pablolagos/go-jsonrpc"
)

const BadRequest = 400

func main() {
	// Define custom options to use a Unix socket and set a custom logger
	options := &go_jsonrpc.Options{
		CGI:    false, // Disable CGI headers
		Logger: log.New(os.Stdout, "JSON-RPC Server: ", log.LstdFlags),
	}

	// Initialize JsRPC server with options
	jsrpc := go_jsonrpc.New(options)

	// Define a global middleware that verifies the presence of a required parameter
	jsrpc.UseGlobalMiddleware(func(ctx *go_jsonrpc.Context) error {
		// Check if the required "authToken" parameter is present
		authToken := ctx.GetParamString("authToken", "")
		if authToken != "secret" {
			ctx.Logger.Println("Unauthorized access attempt.")
			return errors.New("unauthorized")
		}
		return nil
	})

	// Register a command with a specific middleware
	jsrpc.RegisterCommand("multiply", func(ctx *go_jsonrpc.Context) error {
		// Handler logic: Retrieve parameters and return the product
		a := ctx.GetParamFloat("a", 1.0)
		b := ctx.GetParamFloat("b", 1.0)
		result := a * b

		return ctx.JSON(map[string]interface{}{"result": result})
	}, func(ctx *go_jsonrpc.Context) error {
		// Command-specific middleware: Check if parameter "a" is positive
		a := ctx.GetParamFloat("a", 1.0)
		if a <= 0 {
			_ = ctx.Error(BadRequest, errors.New("parameter 'a' must be positive"))
			return errors.New("parameter 'a' must be positive")
		}
		return nil
	})

	// Define the Unix socket path
	socketPath := "/tmp/jsonrpc.sock"
	// Remove socket file if it already exists
	if _, err := os.Stat(socketPath); err == nil {
		os.Remove(socketPath)
	}

	// Start the server on the Unix socket
	log.Printf("Starting server on Unix socket %s", socketPath)
	if err := jsrpc.StartServer(socketPath, true); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
