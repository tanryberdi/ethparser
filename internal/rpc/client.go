package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RPCClient interface defines the methods our RPC client must implement
type RPCClient interface {
	Call(method string, params interface{}) (*JSONRPCResponse, error)
}

type Client struct {
	endpoint   string
	httpClient *http.Client
}

type JSONRPCRequest struct {
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type JSONRPCResponse struct {
	JsonRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result"`
	Error   *JSONRPCError `json:"error"`
	ID      int           `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint:   endpoint,
		httpClient: &http.Client{},
	}
}

func (c *Client) Call(method string, params interface{}) (*JSONRPCResponse, error) {
	request := JSONRPCRequest{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      83,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(c.endpoint, "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var response JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("rpc error: %s", response.Error.Message)
	}

	return &response, nil
}
