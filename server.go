// server.go
package go_jsonrpc

import (
	"encoding/json"
	"io"
)

// ExecuteCommand reads from io.Reader, processes the JSON-RPC request, and writes the response to io.Writer
func (r *JsRPC) ExecuteCommand(reader io.Reader, writer io.Writer) error {
	var rpcRequest JSONRPCRequest

	if err := json.NewDecoder(reader).Decode(&rpcRequest); err != nil {
		return writeError(writer, nil, ParseError, "Invalid JSON", nil, r.cgi)
	}

	// Validate JSON-RPC version
	if rpcRequest.JSONRPC != "2.0" {
		return writeError(writer, rpcRequest.ID, InvalidRequest, "Invalid JSON-RPC version", nil, r.cgi)
	}

	command, exists := r.handlers[rpcRequest.Method]
	if !exists {
		return writeError(writer, rpcRequest.ID, MethodNotFound, "Method not found", nil, r.cgi)
	}

	ctx := &Context{
		Method: rpcRequest.Method,
		Params: rpcRequest.Params,
		writer: writer,
	}

	// Execute global middlewares
	for _, middleware := range r.middlewares {
		if err := middleware(ctx); err != nil {
			return writeError(writer, rpcRequest.ID, InternalError, "Global middleware error", err, r.cgi)
		}
	}

	// Execute command-specific middlewares
	for _, middleware := range command.middlewares {
		if err := middleware(ctx); err != nil {
			return writeError(writer, rpcRequest.ID, InternalError, "Command-specific middleware error", err, r.cgi)
		}
	}

	// Execute handler
	command.handler(ctx)

	// Handle successful response
	return writeResponse(writer, rpcRequest.ID, ctx.Response, r.cgi)
}

// writeHeaders writes the necessary HTTP headers for a CGI response if cgi is true
func writeHeaders(writer io.Writer, cgi bool) error {
	if cgi {
		headers := "Content-Type: application/json\r\n\r\n"
		_, err := writer.Write([]byte(headers))
		return err
	}
	return nil
}

// writeError writes a JSON-RPC 2.0 error response with CGI headers if cgi is true
func writeError(writer io.Writer, id interface{}, code int, message string, data interface{}, cgi bool) error {
	if err := writeHeaders(writer, cgi); err != nil {
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

// writeResponse writes a JSON-RPC 2.0 successful response with CGI headers if cgi is true
func writeResponse(writer io.Writer, id interface{}, result interface{}, cgi bool) error {
	if err := writeHeaders(writer, cgi); err != nil {
		return err
	}
	return json.NewEncoder(writer).Encode(JSONRPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		ID:      id,
	})
}
