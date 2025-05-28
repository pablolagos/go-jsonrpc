# Handler Interceptor in JSON-RPC Server

A **handler interceptor** in the context of the `go-jsonrpc` package serves as a powerful mechanism to modify or halt the request-handling process at an early stage. This functionality allows you to add custom logic that can preprocess incoming requests, perform validations, or even completely take over the request handling by bypassing the primary handler logic.

## Overview of `HandlerInterceptor`

The `HandlerInterceptor` is a user-defined function that operates on incoming requests and responses. It is invoked *before* the JSON-RPC request is processed, giving the developer an opportunity to decide whether or not to proceed with responding.

It is part of the `Options` struct in `go-jsonrpc` and is defined as:

```textmate
HandlerInterceptor func(reader io.Reader, writer io.Writer) (finished bool, err error)
```


### Function Parameters
1. **`reader io.Reader`**: Represents the input stream containing the JSON-RPC request.
2. **`writer io.Writer`**: Represents the output stream where responses or errors may be written.

### Return Values
- **`finished bool`**: Indicates whether the interceptor has completed processing the request. If `true`, the JSON-RPC server will not execute any further handler logic (including middlewares and registered commands).
- **`err error`**: Allows the interceptor to signal an error. The error is logged and the response is sent to the client if applicable.

### Key Features
- **Early Request Validation**: Interceptors can validate input data or authentication tokens early in the request lifecycle.
- **Custom Response Management**: The interceptor can dynamically generate responses or errors without invoking the main handler logic.
- **Preprocessing**: The incoming request stream can be parsed or modified programmatically before reaching the registered handler.

## How to Use `HandlerInterceptor`

The `HandlerInterceptor` can be set in the `Options` struct while initializing the `JsRPC` server. If implemented, every incoming request passes through this function before it reaches other parts of the framework.

### Example 1: Validating Headers
```textmate
options := &go_jsonrpc.Options{
    HandlerInterceptor: func(reader io.Reader, writer io.Writer) (bool, error) {
        // Example: Read the raw JSON request
        requestData, err := io.ReadAll(reader)
        if err != nil {
            return true, fmt.Errorf("failed to read request: %v", err)
        }

        // Validate for custom requirement (e.g., check a token in JSON)
        if !bytes.Contains(requestData, []byte(`"token":"valid_token"`)) {
            // If validation fails, write an error and indicate finished = true
            _, _ = writer.Write([]byte(`{"error": "Invalid token"}`))
            return true, nil
        }

        // Validation passed, allow further handling
        return false, nil
    },
}
```


Here:
- If the token is missing or invalid, a response is immediately written, and the request-handling process stops (`finished = true`).
- If the validation is successful, the request proceeds to other parts of the server (`finished = false`).

### Example 2: Logging and Forcing a Response
```textmate
options := &go_jsonrpc.Options{
    Logger: log.New(os.Stdout, "JSON-RPC Server: ", log.LstdFlags),
    HandlerInterceptor: func(reader io.Reader, writer io.Writer) (bool, error) {
        log.Println("Intercepting request...")

        // Write a fixed response without processing the request
        response := `{"jsonrpc": "2.0", "result": "Intercepted!", "id": null}`
        writer.Write([]byte(response))
        
        // Request handled, no further processing
        log.Println("Request intercepted and response sent.")
        return true, nil
    },
}
```


Here:
- Every request is intercepted, and a pre-defined response (`"Intercepted!"`) is sent without executing any JSON-RPC command or middleware.
- The main handler is bypassed.

### Example 3: Wrapping Streams
Another advanced use case of interceptors could involve modifying or wrapping the `reader` or `writer` streams before the request reaches further stages.

```textmate
options := &go_jsonrpc.Options{
    HandlerInterceptor: func(reader io.Reader, writer io.Writer) (bool, error) {
        // Example: Wrap the writer with logging
        logWriter := io.MultiWriter(writer, os.Stdout)

        // Process the request further with the wrapped writer
        return false, nil // Continue to main handler
    },
}
```


In this example, the output is simultaneously written to both the client and the server's log (using `os.Stdout`).

## When to Use `HandlerInterceptor`
Handler interceptors are particularly beneficial in the following scenarios:
1. **Authentication and Authorization**: Check for tokens or API keys.
2. **Rate-Limiting**: Throttle excessive requests by adding request counters.
3. **Request Validation**: Ensure that the incoming request payload contains required fields or adheres to a schema.
4. **Short-Circuit Responses**: Dynamically respond to certain requests without invoking the main handler logic.
5. **Custom Logging**: Log raw requests and responses explicitly.

## Integration with JSON-RPC Flow

The interceptor is seamlessly integrated into the `go-jsonrpc` request-handling mechanism:

1. The `HandlerInterceptor` is checked in `executeCommandWithData()` before decoding or processing any JSON-RPC request.
2. If the interceptor returns `finished = true`, the request flow terminates there.
3. If `finished = false`, the request proceeds through middlewares, command-specific handlers, and eventually returns a response.

## Conclusion

Handler interceptors provide a highly flexible mechanism for augmenting the JSON-RPC server's behavior. Whether you need to validate requests, log data, enforce policies, or dynamically respond, the interceptor allows you to implement processing at the earliest point in the request lifecycle. Combined with existing middleware and handlers, interceptors make `go-jsonrpc` a powerful framework for building robust JSON-RPC services.