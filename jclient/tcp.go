package jclient

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// TCPClient is a simple JSON-RPC client over raw TCP
type TCPClient struct {
	addr    string
	timeout time.Duration
}

// NewTCPClient creates a new JSON-RPC TCP client
func NewTCPClient(addr string, timeout time.Duration) *TCPClient {
	return &TCPClient{
		addr:    addr,
		timeout: timeout,
	}
}

// Call performs a single JSON-RPC call over TCP and closes the connection
func (c *TCPClient) Call(method string, params interface{}, result interface{}) error {
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
	data = append(data, '\n') // newline as message delimiter

	// connect
	conn, err := net.DialTimeout("tcp", c.addr, c.timeout)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	defer conn.Close()

	// set write deadline
	err = conn.SetWriteDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return fmt.Errorf("set write deadline: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("write: %w", err)
	}

	// read response
	err = conn.SetReadDeadline(time.Now().Add(c.timeout))
	if err != nil {
		return fmt.Errorf("set read deadline: %w", err)
	}

	respBytes, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	// decode
	var rpcResp Response
	if err := json.Unmarshal(respBytes, &rpcResp); err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
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
