// server.go
package go_jsonrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type Logger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// StartServer listens on a TCP port or Unix socket and executes JSON-RPC commands
func (r *JsRPC) StartServer(address string, useUnixSocket bool) error {
	var listener net.Listener
	var err error

	// Select TCP or Unix socket based on the configuration
	if useUnixSocket {
		listener, err = net.Listen("unix", address)
	} else {
		listener, err = net.Listen("tcp", address)
	}
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %v", address, err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			r.logger.Printf("Failed to accept connection: %v", err)
			continue
		}
		go r.handleConnection(conn)
	}
}

// handleConnection reads from a connection, processes the JSON-RPC request, and writes the response
func (r *JsRPC) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Execute command from connection
	if err := r.ExecuteCommand(conn, conn); err != nil {
		r.logger.Printf("Error processing request: %v", err)
	}
}

// ExecuteCommand reads from io.Reader, processes the JSON-RPC request, and writes the response to io.Writer
func (r *JsRPC) ExecuteCommand(reader io.Reader, writer io.Writer) error {
	var rpcRequest JSONRPCRequest

	if err := json.NewDecoder(reader).Decode(&rpcRequest); err != nil {
		r.logger.Printf("Invalid JSON: %v", err)
		return nil
	}

	command, exists := r.handlers[rpcRequest.Method]
	if !exists {
		r.logger.Printf("Method not found: %s", rpcRequest.Method)
		return nil
	}

	ctx := &Context{
		Method: rpcRequest.Method,
		Params: rpcRequest.Params,
		writer: writer,
		Logger: r.logger,
	}

	// Execute global middlewares
	for _, middleware := range r.middlewares {
		if err := middleware(ctx); err != nil {
			// Stop execution if a global middleware returns an error
			return nil
		}
	}

	// Execute command-specific middlewares
	for _, middleware := range command.middlewares {
		if err := middleware(ctx); err != nil {
			// Stop execution if a command-specific middleware returns an error
			return nil
		}
	}

	// Execute handler and log any returned error
	if err := command.handler(ctx); err != nil {
		r.logger.Printf("Handler error: %v", err)
		return nil
	}

	return nil
}

// writeError writes a JSON-RPC 2.0 error response
func writeError(writer io.Writer, id interface{}, code int, message string, data interface{}, cgi bool) error {
	if cgi {
		headers := "Content-Type: application/json\r\n\r\n"
		writer.Write([]byte(headers))
	}
	return json.NewEncoder(writer).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
		ID: id,
	})
}

// writeResponse writes a JSON-RPC 2.0 successful response
func writeResponse(writer io.Writer, id interface{}, result interface{}, cgi bool) error {
	if cgi {
		headers := "Content-Type: application/json\r\n\r\n"
		writer.Write([]byte(headers))
	}
	return json.NewEncoder(writer).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	})
}
