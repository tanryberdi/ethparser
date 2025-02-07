package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ethparser/pkg/types"
)

// MockParser must implement all methods from the Parser interface
type MockParser struct {
	currentBlock int
	subscribers  map[string]bool
	transactions map[string][]types.ParsedTransaction
}

// NewMockParser creates a new mock parser
func NewMockParser() *MockParser {
	return &MockParser{
		currentBlock: 0,
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]types.ParsedTransaction),
	}
}

// Implement Parser interface methods
func (m *MockParser) GetCurrentBlock() int {
	return m.currentBlock
}

func (m *MockParser) Subscribe(address string) bool {
	if m.subscribers[address] {
		return false
	}
	m.subscribers[address] = true
	return true
}

func (m *MockParser) GetTransactions(address string) []types.ParsedTransaction {
	return m.transactions[address]
}

func TestServer(t *testing.T) {
	mockParser := NewMockParser()
	server := NewServer(mockParser)

	t.Run("Subscribe", func(t *testing.T) {
		// Test subscription
		reqBody := SubscribeRequest{Address: "0x123"}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/subscribe", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		server.handleSubscribe(w, req)

		var resp SubscribeResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if !resp.Success {
			t.Errorf("Expected successful subscription")
		}
	})

	t.Run("GetTransactions", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/transactions?address=0x123", nil)
		w := httptest.NewRecorder()

		server.handleGetTransactions(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})

	t.Run("GetCurrentBlock", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/current-block", nil)
		w := httptest.NewRecorder()

		server.handleGetCurrentBlock(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %v", w.Code)
		}
	})
}
