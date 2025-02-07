package parser

import (
	"log"
	"os"
	"testing"

	"ethparser/internal/rpc"
	"ethparser/internal/storage"
)

// MockRPCClient implements the RPCClient interface
type MockRPCClient struct {
	blockNumber int
}

// NewMockRPCClient creates a new mock RPC client
func NewMockRPCClient() *MockRPCClient {
	return &MockRPCClient{
		blockNumber: 1000,
	}
}

// Call implements the RPCClient interface
func (m *MockRPCClient) Call(method string, params interface{}) (*rpc.JSONRPCResponse, error) {
	switch method {
	case "eth_blockNumber":
		return &rpc.JSONRPCResponse{
			JsonRPC: "2.0",
			Result:  "0x3E8", // hex for 1000
			ID:      1,
		}, nil
	case "eth_getBlockByNumber":
		// Create a mock transaction that matches our subscribed address
		return &rpc.JSONRPCResponse{
			JsonRPC: "2.0",
			Result: map[string]interface{}{
				"number":    "0x3E8",
				"hash":      "0x123",
				"timestamp": "0x60c88c32",
				"transactions": []map[string]interface{}{
					{
						"hash":        "0xabc",
						"from":        "0x123456",
						"to":          "0xdac17f958d2ee523a2206206994597c13d831ec7", // This matches our test address
						"value":       "0x0",
						"blockNumber": "0x3E8",
					},
				},
			},
			ID: 1,
		}, nil
	}
	return nil, nil
}

// Helper function to create a new parser instance for each test
func createTestParser() *EthParser {
	logger := log.New(os.Stdout, "test: ", log.LstdFlags)
	mockClient := NewMockRPCClient()
	storage := storage.NewMemoryStorage(logger)

	return &EthParser{
		client:  mockClient,
		storage: storage,
		logger:  logger,
	}
}

func TestParser(t *testing.T) {
	// Test Subscribe (separate test with its own parser instance)
	t.Run("Subscribe", func(t *testing.T) {
		parser := createTestParser()
		address := "0xdac17f958d2ee523a2206206994597c13d831ec7"

		// First subscription should succeed
		if !parser.Subscribe(address) {
			t.Error("First subscription should succeed")
		}

		// Second subscription should fail
		if parser.Subscribe(address) {
			t.Error("Second subscription should fail")
		}
	})

	// Test GetCurrentBlock (separate test with its own parser instance)
	t.Run("GetCurrentBlock", func(t *testing.T) {
		parser := createTestParser()
		block := parser.GetCurrentBlock()
		if block < 0 {
			t.Errorf("Block number should not be negative, got %d", block)
		}
	})

	// Test GetTransactions (separate test with its own parser instance)
	t.Run("GetTransactions", func(t *testing.T) {
		parser := createTestParser()
		address := "0xdac17f958d2ee523a2206206994597c13d831ec7"

		// First, subscribe to the address
		if !parser.Subscribe(address) {
			t.Error("Failed to subscribe address")
		}

		// Start the parser
		err := parser.Start()
		if err != nil {
			t.Errorf("Failed to start parser: %v", err)
		}

		// Manually trigger block parsing
		err = parser.parseBlock(1000)
		if err != nil {
			t.Errorf("Failed to parse block: %v", err)
		}

		// Get transactions
		transactions := parser.GetTransactions(address)
		if len(transactions) == 0 {
			t.Error("Expected to find transactions")
		} else {
			// Verify transaction details
			tx := transactions[0]
			if tx.To != address {
				t.Errorf("Expected transaction to address %s, got %s", address, tx.To)
			}
		}
	})
}
