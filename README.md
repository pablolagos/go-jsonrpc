![Go Version](https://img.shields.io/github/go-mod/go-version/pablolagos/go-jsonrpc)
![License](https://img.shields.io/github/license/pablolagos/go-jsonrpc)
![Issues](https://img.shields.io/github/issues/pablolagos/go-jsonrpc)
![Code Size](https://img.shields.io/github/languages/code-size/pablolagos/go-jsonrpc)
![Last Commit](https://img.shields.io/github/last-commit/pablolagos/go-jsonrpc)
![Go Report Card](https://goreportcard.com/badge/github.com/pablolagos/go-jsonrpc)

![Stars](https://img.shields.io/github/stars/pablolagos/go-jsonrpc?style=social)
![Stars](https://img.shields.io/github/stars/pablolagos/go-jsonrpc?style=social)
![Forks](https://img.shields.io/github/forks/pablolagos/go-jsonrpc?style=social)
[![Documentation](https://img.shields.io/badge/docs-available-brightgreen)](https://pkg.go.dev/github.com/pablolagos/go-jsonrpc)


# Go JSON-RPC

A versatile and lightweight JSON-RPC 2.0 server implementation in Go, designed to handle JSON-RPC requests over TCP, Unix sockets, CGI, and HTTP. It provides a full set of helper functions to register commands, handle middlewares, and work with JSON-RPC data in a flexible and customizable way.

## Features

- **Multiple transport protocols**: TCP, Unix sockets, CGI, and HTTP.
- **Middleware support**: Easily add global or command-specific middlewares.
- **Helper functions**: Simplify request handling, parameter retrieval, and response management.
- **Flexible configuration**: Customize logging, error handling, and server options.

## Installation

To install this package, use:

```bash
go get github.com/pablolagos/go-jsonrpc
```

## Usage

### Basic Server Setup (TCP and Unix Sockets)

#### 1. TCP Server Example

```go
package main

import (
    "errors"
    "log"
    "os"
    "github.com/pablolagos/go-jsonrpc"
)

func main() {
    options := &go_jsonrpc.Options{
        CGI:    false,
        Logger: log.New(os.Stdout, "JSON-RPC TCP Server: ", log.LstdFlags),
    }

    jsrpc := go_jsonrpc.New(options)

    // Register a command
    jsrpc.RegisterCommand("add", func(ctx *go_jsonrpc.Context) error {
        a := ctx.GetParamFloat("a", 0.0)
        b := ctx.GetParamFloat("b", 0.0)
        result := a + b
        return ctx.JSON(result)
    })

    // Start the server on a TCP port
    address := ":12345"
    log.Printf("Starting server on %s\n", address)
    if err := jsrpc.StartServer(address, false); err != nil {
        log.Fatalf("Failed to start server: %v\n", err)
    }
}
```

#### 2. Unix Socket Example

```go
package main

import (
    "errors"
    "log"
    "os"
    "github.com/pablolagos/go-jsonrpc"
)

func main() {
    options := &go_jsonrpc.Options{
        CGI:    false,
        Logger: log.New(os.Stdout, "JSON-RPC Unix Socket Server: ", log.LstdFlags),
    }

    jsrpc := go_jsonrpc.New(options)

    jsrpc.UseGlobalMiddleware(func(ctx *go_jsonrpc.Context) error {
        authToken := ctx.GetParamString("authToken", "")
        if authToken != "secret" {
            return errors.New("unauthorized access")
        }
        return nil
    })

    jsrpc.RegisterCommand("multiply", func(ctx *go_jsonrpc.Context) error {
        a := ctx.GetParamFloat("a", 1.0)
        b := ctx.GetParamFloat("b", 1.0)
        result := a * b
        ctx.JSON(result)
        return nil
    })

    socketPath := "/tmp/jsonrpc.sock"
    if _, err := os.Stat(socketPath); err == nil {
        os.Remove(socketPath)
    }

    log.Printf("Starting server on Unix socket %s\n", socketPath)
    if err := jsrpc.StartServer(socketPath, true); err != nil {
        log.Fatalf("Failed to start server: %v\n", err)
    }
}
```

### CGI Example

When running in a CGI environment, the library can automatically handle JSON-RPC requests by reading from `os.Stdin` and writing responses to `os.Stdout`.

```go
package main

import (
    "log"
    "os"
    "github.com/pablolagos/go-jsonrpc"
)

func main() {
    options := &go_jsonrpc.Options{
        CGI:    true, // Enables CGI headers
        Logger: log.New(os.Stdout, "JSON-RPC CGI Server: ", log.LstdFlags),
    }

    jsrpc := go_jsonrpc.New(options)

    jsrpc.RegisterCommand("subtract", func(ctx *go_jsonrpc.Context) error {
        a := ctx.GetParamFloat("a", 0.0)
        b := ctx.GetParamFloat("b", 0.0)
        result := a - b
        return ctx.JSON(result)
    })

    if err := jsrpc.ExecuteCommand(os.Stdin, os.Stdout); err != nil {
        log.Fatalf("Error processing CGI request: %v", err)
    }
}
```

### HTTP Server Example

To run the server over HTTP, use Go's `net/http` package alongside the JSON-RPC library.

```go
package main

import (
    "log"
    "net/http"
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
            return ctx.ErrorString(400,"division by zero")
        }
        result := a / b
        return ctx.JSON(result)
    })

    http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
        jsrpc.ExecuteCommand(r.Body, w)
    })

    log.Println("Starting HTTP server on :8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatalf("Failed to start HTTP server: %v", err)
    }
}
```

## Middleware Example

Middlewares can be applied globally or specifically for individual commands. Note: If a middleware encounters an error, it is responsible for handling client responses as needed.

```go
package main

import (
    "errors"
    "log"
    "os"
    "github.com/pablolagos/go-jsonrpc"
)

func main() {
    options := &go_jsonrpc.Options{
        Logger: log.New(os.Stdout, "JSON-RPC Server with Middleware: ", log.LstdFlags),
    }

    jsrpc := go_jsonrpc.New(options)

    jsrpc.UseGlobalMiddleware(func(ctx *go_jsonrpc.Context) error {
        token := ctx.GetParamString("token", "")
        if token != "valid_token" {
            ctx.ErrorString(412, "unauthorized access")
            return errors.New("unauthorized access") // Return an error to stop command execution
        }
        return nil
    })

    jsrpc.RegisterCommand("echo", func(ctx *go_jsonrpc.Context) error {
        message := ctx.GetParamString("message", "")
        return ctx.JSON(message)
    })

    address := ":12345"
    log.Printf("Starting server on %s\n", address)
    if err := jsrpc.StartServer(address, false); err != nil {
        log.Fatalf("Failed to start server: %v\n", err)
    }
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
