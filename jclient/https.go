package jclient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Request represents a JSON-RPC request
type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      int         `json:"id"`
}

// Response represents a JSON-RPC response
type Response struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result,omitempty"`
	Error   *RPCError        `json:"error,omitempty"`
	ID      int              `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPClient is a simple JSON-RPC client over HTTP/HTTPS
type HTTPClient struct {
	endpoint string
	client   *http.Client
}

type HTTPClientOptions struct {
	Insecure bool          // Allow self-signed certificates
	Timeout  time.Duration // HTTP client timeout
}

// NewHTTPClient creates a new JSON-RPC client with HTTPS support
func NewHTTPClient(endpoint string, opts HTTPClientOptions) *HTTPClient {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opts.Insecure}, // allow self-signed certs if insecure==true
	}
	return &HTTPClient{
		endpoint: endpoint,
		client: &http.Client{
			Transport: tr,
			Timeout:   opts.Timeout,
		},
	}
}

// Call performs a JSON-RPC call and decodes into result
func (c *HTTPClient) Call(method string, params interface{}, result interface{}) error {
	req := Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	var rpcResp Response
	if err := json.NewDecoder(resp.Body).Decode(&rpcResp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if rpcResp.Error != nil {
		return fmt.Errorf("rpc error %d: %s", rpcResp.Error.Code, rpcResp.Error.Message)
	}

	if rpcResp.Result != nil && result != nil {
		if err := json.Unmarshal(*rpcResp.Result, result); err != nil {
			return fmt.Errorf("unmarshal result: %w", err)
		}
	}

	return nil
}
