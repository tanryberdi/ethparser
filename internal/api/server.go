package api

import (
	"encoding/json"
	"log"
	"net/http"

	"ethparser/pkg/types"
)

type Server struct {
	parser types.Parser
}

func NewServer(parser types.Parser) *Server {
	return &Server{
		parser: parser,
	}
}

type SubscribeRequest struct {
	Address string `json:"address"`
}

type SubscribeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type GetTransactionsResponse struct {
	Address      string                    `json:"address"`
	Transactions []types.ParsedTransaction `json:"transactions"`
}

func (s *Server) RegisterRoutes() {
	http.HandleFunc("/subscribe", s.handleSubscribe)
	http.HandleFunc("/transactions", s.handleGetTransactions)
	http.HandleFunc("/current-block", s.handleGetCurrentBlock)
}

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SubscribeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Subscribing to address: %s", req.Address)
	success := s.parser.Subscribe(req.Address)
	log.Printf("Subscription result for %s: %v", req.Address, success)
	resp := SubscribeResponse{
		Success: success,
		Message: getSubscribeMessage(success),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleGetTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	address := r.URL.Query().Get("address")
	if address == "" {
		log.Printf("Missing address parameter")
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	log.Printf("Getting transactions for address: %s", address)
	transactions := s.parser.GetTransactions(address)
	log.Printf("Found %d transactions for address %s", len(transactions), address)

	resp := GetTransactionsResponse{
		Address:      address,
		Transactions: transactions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleGetCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	currentBlock := s.parser.GetCurrentBlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"currentBlock": currentBlock})
}

func getSubscribeMessage(success bool) string {
	if success {
		return "Address subscribed successfully"
	}
	return "Address already subscribed"
}
