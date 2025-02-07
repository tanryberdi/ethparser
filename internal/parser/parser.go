package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"ethparser/internal/rpc"
	"ethparser/internal/storage"
	"ethparser/pkg/types"
)

type EthParser struct {
	client  rpc.RPCClient
	storage *storage.MemoryStorage
	logger  *log.Logger
}

func NewEthParser(endpoint string, logger *log.Logger) *EthParser {
	return &EthParser{
		client:  rpc.NewClient(endpoint),
		storage: storage.NewMemoryStorage(logger),
		logger:  logger,
	}
}

func (p *EthParser) GetCurrentBlock() int {
	return p.storage.GetCurrentBlock()
}

func (p *EthParser) Subscribe(address string) bool {
	return p.storage.Subscribe(address)
}

func (p *EthParser) GetTransactions(address string) []types.ParsedTransaction {
	return p.storage.GetTransactions(address)
}

func (p *EthParser) Start() error {
	// Get latest block number first
	resp, err := p.client.Call("eth_blockNumber", []interface{}{})
	if err != nil {
		return fmt.Errorf("failed to get latest block: %w", err)
	}

	latestBlock := hexToInt(resp.Result.(string))
	p.logger.Printf("Starting from latest block: %d", latestBlock)

	// Set as our starting point
	p.storage.SetCurrentBlock(latestBlock)

	// Start blockchain parsing in a separate goroutine
	go p.parseBlocks()
	return nil
}

func (p *EthParser) parseBlocks() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	p.logger.Printf("Starting block parser...")

	for range ticker.C {
		currentBlock := p.GetCurrentBlock()
		p.logger.Printf("Current block: %d", currentBlock)

		// Get latest block number
		resp, err := p.client.Call("eth_blockNumber", []interface{}{})
		if err != nil {
			p.logger.Printf("Failed to get latest block: %v", err)
			continue
		}

		latestBlock := hexToInt(resp.Result.(string))
		p.logger.Printf("Latest block from network: %d", latestBlock)

		if currentBlock < latestBlock {
			p.logger.Printf("Processing blocks from %d to %d", currentBlock+1, latestBlock)

			// Parse new blocks
			for blockNum := currentBlock + 1; blockNum <= latestBlock; blockNum++ {
				p.logger.Printf("Parsing block %d", blockNum)
				if err := p.parseBlock(blockNum); err != nil {
					p.logger.Printf("Failed to parse block %d: %v", blockNum, err)
					continue
				}
				p.storage.SetCurrentBlock(blockNum)
			}
		} else {
			p.logger.Printf("No new blocks to process")
		}
	}
}

// Block represents an Ethereum block structure
type Block struct {
	Number       string              `json:"number"`
	Hash         string              `json:"hash"`
	Timestamp    string              `json:"timestamp"`
	Transactions []types.Transaction `json:"transactions"`
}

// Transaction represents an Ethereum transaction structure
type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber string `json:"blockNumber"`
}

func (p *EthParser) parseBlock(blockNum int) error {
	p.logger.Printf("Starting to parse block %d", blockNum)

	blockHex := fmt.Sprintf("0x%x", blockNum)

	resp, err := p.client.Call("eth_getBlockByNumber", []interface{}{blockHex, true})
	if err != nil {
		return fmt.Errorf("failed to get block %d: %w", blockNum, err)
	}

	var block Block
	blockData, err := json.Marshal(resp.Result)
	if err != nil {
		return fmt.Errorf("failed to marshal block data: %w", err)
	}

	if err := json.Unmarshal(blockData, &block); err != nil {
		return fmt.Errorf("failed to unmarshal block data: %w", err)
	}

	p.logger.Printf("Processing %d transactions in block %d", len(block.Transactions), blockNum)

	transactionsFound := 0
	for _, tx := range block.Transactions {
		// Debug logging
		p.logger.Printf("Checking transaction: From=%s, To=%s", tx.From, tx.To)

		if p.storage.IsSubscribed(tx.From) || p.storage.IsSubscribed(tx.To) {
			p.logger.Printf("Found relevant transaction in block %d: %s", blockNum, tx.Hash)

			// Convert to ParsedTransaction
			parsedTx := types.ParsedTransaction{
				Hash:        tx.Hash,
				From:        tx.From,
				To:          tx.To,
				Value:       tx.Value,
				BlockNumber: int64(blockNum),
				Timestamp:   int64(hexToInt(block.Timestamp)),
			}

			p.storage.AddTransaction(parsedTx)
			transactionsFound++
		}
	}

	p.logger.Printf("Found %d relevant transactions in block %d", transactionsFound, blockNum)
	return nil
}

func hexToInt(hex string) int {
	// Remove "0x" prefix if present
	if len(hex) >= 2 && hex[0:2] == "0x" {
		hex = hex[2:]
	}

	// Convert hex string to big.Int
	n := new(big.Int)
	n.SetString(hex, 16)

	return int(n.Int64())
}
