package main

import (
	"log"
	"net/http"
	"os"

	"github.com/pablolagos/go-jsonrpc"
)

func main() {
	options := &go_jsonrpc.Options{
		CGI:    false,
		Logger: log.New(os.Stdout, "JSON-RPC HTTP Server: ", log.LstdFlags),
	}

	jsrpc := go_jsonrpc.New(options)

	jsrpc.RegisterCommand("divide", func(ctx *go_jsonrpc.Context) error {
		a := ctx.GetParamFloat("a", 1.0)
		b := ctx.GetParamFloat("b", 1.0)
		if b == 0 {
			return ctx.ErrorString(go_jsonrpc.InternalError, "division by zero")
		}
		result := a / b
		ctx.JSON(map[string]interface{}{"result": result})
		return nil
	})

	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		// Handle preflight OPTIONS request for CORS (optional)
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.WriteHeader(http.StatusNoContent) // No content for preflight response
			return
		}

		// Add CORS headers to the actual response (optional)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json") // Set your preferred content type

		jsrpc.ExecuteCommand(r.Body, w)
	})

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
