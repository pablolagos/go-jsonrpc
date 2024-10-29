package cgi

import (
	"fmt"
	"github.com/pablolagos/go-jsonrpc"
	"os"
)

// UserInfo represents a struct for binding JSON parameters
type UserInfo struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Active bool     `json:"active"`
	Roles  []string `json:"roles"`
}

func main() {
	// Create an instance of JsRPC with the CGI flag set to true to write CGI headers
	jsrpc := go_jsonrpc.New(true)

	// Global middleware to add a request ID to the context
	jsrpc.UseGlobalMiddleware(func(ctx *go_jsonrpc.Context) error {
		// Store a unique request ID in the context
		ctx.SetData("request_id", 1001)
		fmt.Println("Global middleware executed")
		return nil
	})

	// Register a "sum" command that adds two numbers using GetParamFloat
	jsrpc.RegisterCommand("sum", func(ctx *go_jsonrpc.Context) {
		// Retrieve individual parameters using GetParamFloat
		a := ctx.GetParamFloat("a", 0.0)
		b := ctx.GetParamFloat("b", 0.0)

		// Retrieve the request ID from the context
		requestID := ctx.GetData("request_id")

		// Perform the operation and return the response
		ctx.JSON(map[string]interface{}{
			"request_id": requestID,
			"result":     a + b,
		})
	})

	// Register a "getUserInfo" command that binds JSON parameters to a struct
	jsrpc.RegisterCommand("getUserInfo", func(ctx *go_jsonrpc.Context) {
		// Use Bind to parse parameters into a UserInfo struct
		var userInfo UserInfo
		if err := ctx.Bind(&userInfo); err != nil {
			ctx.Error(fmt.Errorf("failed to bind parameters: %v", err))
			return
		}

		// Return the bound user information as a JSON response
		ctx.JSON(map[string]interface{}{
			"user_id": userInfo.ID,
			"name":    userInfo.Name,
			"active":  userInfo.Active,
			"roles":   userInfo.Roles,
		})
	})

	// Execute the command using os.Stdin and os.Stdout as input and output (for CGI)
	if err := jsrpc.ExecuteCommand(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
