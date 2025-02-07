package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	client := NewClient("http://localhost:8080")

	// Test addresses (using known active addresses)
	addresses := []string{
		"0x8ad599c3A0ff1De082011EFDDc58f1908eb6e6D8",
		"0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D",
	}

	// Run tests
	for _, address := range addresses {
		fmt.Printf("\n=== Testing with address: %s ===\n", address)

		// 1. Subscribe to address
		fmt.Println("\nTesting Subscribe:")
		resp, err := client.Subscribe(address)
		if err != nil {
			log.Printf("Subscribe error: %v", err)
		} else {
			fmt.Printf("Subscribe response: %+v\n", resp)
		}

		// 2. Get initial block
		fmt.Println("\nGetting initial block:")
		initialBlock, err := client.GetCurrentBlock()
		if err != nil {
			log.Printf("Get current block error: %v", err)
		} else {
			fmt.Printf("Current block: %d\n", initialBlock)
		}

		// 3. Wait for some blocks to be processed
		fmt.Println("\nWaiting 30 seconds for blocks to be processed...")
		time.Sleep(30 * time.Second)

		// 4. Get transactions
		fmt.Println("\nGetting transactions:")
		txResp, err := client.GetTransactions(address)
		if err != nil {
			log.Printf("Get transactions error: %v", err)
		} else {
			fmt.Printf("Found %d transactions\n", len(txResp.Transactions))
			for i, tx := range txResp.Transactions {
				fmt.Printf("\nTransaction %d:\n", i+1)
				fmt.Printf("  Hash: %s\n", tx.Hash)
				fmt.Printf("  From: %s\n", tx.From)
				fmt.Printf("  To: %s\n", tx.To)
				fmt.Printf("  Value: %s\n", tx.Value)
				fmt.Printf("  Block: %d\n", tx.BlockNumber)
				fmt.Printf("  Time: %s\n", time.Unix(tx.Timestamp, 0))
			}
		}

		// 5. Get final block
		fmt.Println("\nGetting final block:")
		finalBlock, err := client.GetCurrentBlock()
		if err != nil {
			log.Printf("Get current block error: %v", err)
		} else {
			fmt.Printf("Final block: %d\n", finalBlock)
			fmt.Printf("Blocks processed: %d\n", finalBlock-initialBlock)
		}
	
	}

}
