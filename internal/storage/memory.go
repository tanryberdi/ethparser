package storage

import (
	"log"
	"strings"
	"sync"

	"ethparser/pkg/types"
)

type MemoryStorage struct {
	mu           sync.RWMutex
	subscribers  map[string]bool
	transactions map[string][]types.ParsedTransaction
	currentBlock int
	logger       *log.Logger
}

func NewMemoryStorage(logger *log.Logger) *MemoryStorage {
	return &MemoryStorage{
		subscribers:  make(map[string]bool),
		transactions: make(map[string][]types.ParsedTransaction),
		currentBlock: 0,
		logger:       logger,
	}
}

func (s *MemoryStorage) IsSubscribed(address string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Debug: Print all current subscribers
	s.logger.Printf("Current subscribers: %+v", s.subscribers)

	subscribed := s.subscribers[address]
	s.logger.Printf("Checking subscription for address %s: %v", address, subscribed)
	return subscribed
}

func (s *MemoryStorage) Subscribe(address string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Convert address to lowercase for consistent comparison
	address = strings.ToLower(address)

	s.logger.Printf("Attempting to subscribe address: %s", address)
	s.logger.Printf("Current subscribers before: %+v", s.subscribers)

	if _, exists := s.subscribers[address]; exists {
		s.logger.Printf("Address %s is already subscribed", address)
		return false
	}

	s.subscribers[address] = true
	s.logger.Printf("Current subscribers after: %+v", s.subscribers)
	s.logger.Printf("Successfully subscribed address: %s", address)
	return true
}

//func (s *MemoryStorage) GetTransactions(address string) []types.Transaction {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	s.logger.Printf("Getting transactions for address: %s", address)
//	txs := s.transactions[address]
//	s.logger.Printf("Found %d transactions for address %s", len(txs), address)
//	return txs
//}
//
//func (s *MemoryStorage) AddTransaction(tx types.Transaction) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	s.logger.Printf("Adding transaction: Hash=%s, From=%s, To=%s", tx.Hash, tx.From, tx.To)
//
//	if tx.From != "" && s.subscribers[tx.From] {
//		s.transactions[tx.From] = append(s.transactions[tx.From], tx)
//		s.logger.Printf("Added outgoing transaction for %s", tx.From)
//	}
//
//	if tx.To != "" && s.subscribers[tx.To] {
//		s.transactions[tx.To] = append(s.transactions[tx.To], tx)
//		s.logger.Printf("Added incoming transaction for %s", tx.To)
//	}
//}

func (s *MemoryStorage) AddTransaction(tx types.ParsedTransaction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Printf("Adding transaction: Hash=%s, From=%s, To=%s", tx.Hash, tx.From, tx.To)

	if tx.From != "" && s.subscribers[tx.From] {
		s.transactions[tx.From] = append(s.transactions[tx.From], tx)
		s.logger.Printf("Added outgoing transaction for %s", tx.From)
	}

	if tx.To != "" && s.subscribers[tx.To] {
		s.transactions[tx.To] = append(s.transactions[tx.To], tx)
		s.logger.Printf("Added incoming transaction for %s", tx.To)
	}
}

func (s *MemoryStorage) GetTransactions(address string) []types.ParsedTransaction {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logger.Printf("Getting transactions for address: %s", address)
	txs := s.transactions[address]
	s.logger.Printf("Found %d transactions for address %s", len(txs), address)
	return txs
}

func (s *MemoryStorage) SetCurrentBlock(block int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.logger.Printf("Updating current block from %d to %d", s.currentBlock, block)
	s.currentBlock = block
}

func (s *MemoryStorage) GetCurrentBlock() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.logger.Printf("Getting current block: %d", s.currentBlock)
	return s.currentBlock
}
