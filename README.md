# go_jsonrpc

**go_jsonrpc** is a lightweight and flexible library for building JSON-RPC 2.0 APIs in Go. Designed with simplicity and modularity in mind, it provides everything needed to set up and manage commands, middleware, and CGI headers for diverse environments. It also supports parameter extraction, type-safe binding, and shared data between middlewares and handlers.

## Features

- **JSON-RPC 2.0 Compatibility**: Full support for JSON-RPC 2.0 requests and responses.
- **Middleware Support**: Global and command-specific middlewares for flexible request handling.
- **Parameter Extraction**: Easily retrieve parameters with type safety.
- **Struct Binding**: Use `Bind` to convert JSON parameters to custom Go structs.
- **CGI Support**: Optional CGI headers for running in environments that require HTTP-like responses.
- **Shared Context Data**: Share information between middlewares and handlers using `SetData` and `GetData`.

## Installation

To install **go_jsonrpc**, you can use `go get`:

```sh
go get github.com/pablolagos/go_jsonrpc
```

## Usage

Here's a quick example to get you started. This example demonstrates how to:
1. Create an instance of `JsRPC`.
2. Set up global and command-specific middleware.
3. Register commands that use parameter extraction, `Bind`, and shared context data.
4. Execute commands and produce JSON-RPC responses.

### Example

```go
package main

import (
    "fmt"
    "os"
    "github.com/pablolagos/go_jsonrpc"
)

// UserInfo represents a struct for binding JSON parameters
type UserInfo struct {
    ID     int      `json:"id"`
    Name   string   `json:"name"`
    Active bool     `json:"active"`
    Roles  []string `json:"roles"`
}

func main() {
    // Create an instance of JsRPC with CGI headers enabled
    jsrpc := go_jsonrpc.New(true)

    // Global middleware to add a request ID to the context
    jsrpc.UseGlobalMiddleware(func(ctx *go_jsonrpc.Context) error {
        ctx.SetData("request_id", 1001)
        fmt.Println("Global middleware executed")
        return nil
    })

    // Register a "sum" command that uses GetParamFloat
    jsrpc.RegisterCommand("sum", func(ctx *go_jsonrpc.Context) {
        a := ctx.GetParamFloat("a", 0.0)
        b := ctx.GetParamFloat("b", 0.0)
        requestID := ctx.GetData("request_id")

        ctx.JSON(map[string]interface{}{
            "request_id": requestID,
            "result":     a + b,
        })
    })

    // Register a "getUserInfo" command that binds JSON parameters to a struct
    jsrpc.RegisterCommand("getUserInfo", func(ctx *go_jsonrpc.Context) {
        var userInfo UserInfo
        if err := ctx.Bind(&userInfo); err != nil {
            ctx.Error(fmt.Errorf("failed to bind parameters: %v", err))
            return
        }

        ctx.JSON(map[string]interface{}{
            "user_id": userInfo.ID,
            "name":    userInfo.Name,
            "active":  userInfo.Active,
            "roles":   userInfo.Roles,
        })
    })

    // Execute the command using os.Stdin and os.Stdout (for CGI environments)
    if err := jsrpc.ExecuteCommand(os.Stdin, os.Stdout); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    }
}
```

### Example JSON-RPC Requests

#### Request for `sum` Command

```json
{
  "jsonrpc": "2.0",
  "method": "sum",
  "params": { "a": 10.5, "b": 5.5 },
  "id": 1
}
```

#### Response

```json
{
  "jsonrpc": "2.0",
  "result": {
    "request_id": 1001,
    "result": 16.0
  },
  "id": 1
}
```

#### Request for `getUserInfo` Command

```json
{
  "jsonrpc": "2.0",
  "method": "getUserInfo",
  "params": {
    "id": 42,
    "name": "Alice",
    "active": true,
    "roles": ["admin", "user"]
  },
  "id": 2
}
```

#### Response

```json
{
  "jsonrpc": "2.0",
  "result": {
    "user_id": 42,
    "name": "Alice",
    "active": true,
    "roles": ["admin", "user"]
  },
  "id": 2
}
```

## API Reference

#### New(cgi bool) *JsRPC
Creates a new `JsRPC` instance. If `cgi` is set to `true`, CGI headers will be included in responses.

#### RegisterCommand(commandName string, handler HandlerFunc, middlewares ...MiddlewareFunc)
Registers a new command with optional specific middlewares.

#### UseGlobalMiddleware(middleware MiddlewareFunc)
Adds a global middleware that will apply to all commands.

#### ExecuteCommand(reader io.Reader, writer io.Writer) error
Reads a JSON-RPC request from `reader`, processes it, and writes the response to `writer`. It includes CGI headers if `cgi` is set to `true`.

## Contributing

Feel free to open issues or pull requests if you have any improvements, suggestions, or bug fixes! Contributions are always welcome.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

