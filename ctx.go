package go_jsonrpc

import (
	"encoding/json"
	"io"
	"log"
)

type Context struct {
	Method   string         // The method being executed
	Params   any            // Params can be either an array or a map
	Response interface{}    // The response to be sent
	writer   io.Writer      // Writer for the response
	data     map[string]any // To store shared data between middleware and handlers
	Logger   *log.Logger    // Logger available for handlers and middlewares
}

// JSON writes a JSON response to the writer
func (ctx *Context) JSON(data interface{}) error {
	ctx.Response = data
	return json.NewEncoder(ctx.writer).Encode(ctx.Response)
}

// Error writes an error response to the writer
func (ctx *Context) Error(err error) error {
	response := map[string]string{
		"error": err.Error(),
	}
	return ctx.JSON(response)
}

// Bind binds the params to the provided destination struct
func (ctx *Context) Bind(dest interface{}) error {
	bytes, err := json.Marshal(ctx.Params)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, dest)
}

// SetData stores a value in the context that can be shared across middlewares and handlers
func (ctx *Context) SetData(name string, value any) {
	if ctx.data == nil {
		ctx.data = make(map[string]any)
	}
	ctx.data[name] = value
}

// GetData retrieves a value stored in the context by name
func (ctx *Context) GetData(name string) any {
	if ctx.data != nil {
		return ctx.data[name]
	}
	return nil
}

// GetParamInt retrieves an integer parameter by name, or returns a default value if not found
func (ctx *Context) GetParamInt(name string, defaultValue int) int {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if intVal, ok := val.(float64); ok { // JSON numbers are float64 in Go
				return int(intVal)
			}
		}
	}
	return defaultValue
}

// GetParamFloat retrieves a float64 parameter by name, or returns a default value if not found
func (ctx *Context) GetParamFloat(name string, defaultValue float64) float64 {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if floatVal, ok := val.(float64); ok {
				return floatVal
			}
		}
	}
	return defaultValue
}

// GetParamString retrieves a string parameter by name, or returns a default value if not found
func (ctx *Context) GetParamString(name string, defaultValue string) string {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if strVal, ok := val.(string); ok {
				return strVal
			}
		}
	}
	return defaultValue
}

// GetParamBool retrieves a boolean parameter by name, or returns a default value if not found
func (ctx *Context) GetParamBool(name string, defaultValue bool) bool {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if boolVal, ok := val.(bool); ok {
				return boolVal
			}
		}
	}
	return defaultValue
}

// GetParamIntArray retrieves a slice of integers by name, or returns a default empty slice if not found
func (ctx *Context) GetParamIntArray(name string) []int {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if arrayVal, ok := val.([]interface{}); ok {
				result := make([]int, 0, len(arrayVal))
				for _, item := range arrayVal {
					if intVal, ok := item.(float64); ok { // JSON numbers are float64
						result = append(result, int(intVal))
					}
				}
				return result
			}
		}
	}
	return []int{}
}

// GetParamFloatArray retrieves a slice of float64 by name, or returns a default empty slice if not found
func (ctx *Context) GetParamFloatArray(name string) []float64 {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if arrayVal, ok := val.([]interface{}); ok {
				result := make([]float64, 0, len(arrayVal))
				for _, item := range arrayVal {
					if floatVal, ok := item.(float64); ok {
						result = append(result, floatVal)
					}
				}
				return result
			}
		}
	}
	return []float64{}
}

// GetParamStringArray retrieves a slice of strings by name, or returns a default empty slice if not found
func (ctx *Context) GetParamStringArray(name string) []string {
	if paramsMap, ok := ctx.Params.(map[string]interface{}); ok {
		if val, found := paramsMap[name]; found {
			if arrayVal, ok := val.([]interface{}); ok {
				result := make([]string, 0, len(arrayVal))
				for _, item := range arrayVal {
					if strVal, ok := item.(string); ok {
						result = append(result, strVal)
					}
				}
				return result
			}
		}
	}
	return []string{}
}
