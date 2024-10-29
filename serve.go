package go_jsonrpc

import (
	"encoding/json"
	"io"
)

// ExecuteCommand reads from the provided io.Reader, processes the JSON-RPC request, and writes the response to io.Writer
func (r *JsRPC) ExecuteCommand(reader io.Reader, writer io.Writer) error {
	var rpcRequest JSONRPCRequest

	if err := json.NewDecoder(reader).Decode(&rpcRequest); err != nil {
		return writeError(writer, nil, ParseError, "Invalid JSON", nil)
	}

	// Validate JSON-RPC version
	if rpcRequest.JSONRPC != "2.0" {
		return writeError(writer, rpcRequest.ID, InvalidRequest, "Invalid JSON-RPC version", nil)
	}

	command, exists := r.handlers[rpcRequest.Method]
	if !exists {
		return writeError(writer, rpcRequest.ID, MethodNotFound, "Method not found", nil)
	}

	ctx := &Context{
		Method: rpcRequest.Method,
		Params: rpcRequest.Params,
		writer: writer,
	}

	// Run global middlewares
	for _, middleware := range r.middlewares {
		if err := middleware(ctx); err != nil {
			return writeError(writer, rpcRequest.ID, InternalError, "Global middleware error", err)
		}
	}

	// Run command-specific middlewares
	for _, middleware := range command.Middlewares {
		if err := middleware(ctx); err != nil {
			return writeError(writer, rpcRequest.ID, InternalError, "Command-specific middleware error", err)
		}
	}

	// Execute handler
	command.Handler(ctx)

	// Handle successful response
	return writeResponse(writer, rpcRequest.ID, ctx.Response)
}

// writeHeaders writes the necessary HTTP headers for a CGI response
func writeHeaders(writer io.Writer) error {
	headers := "Content-Type: application/json\r\n\r\n"
	_, err := writer.Write([]byte(headers))
	return err
}

// writeError writes a JSON-RPC 2.0 error response with CGI headers
func writeError(writer io.Writer, id interface{}, code int, message string, data interface{}) error {
	if err := writeHeaders(writer); err != nil {
		return err
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

// writeResponse writes a JSON-RPC 2.0 successful response with CGI headers
func writeResponse(writer io.Writer, id interface{}, result interface{}) error {
	if err := writeHeaders(writer); err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	})
}
