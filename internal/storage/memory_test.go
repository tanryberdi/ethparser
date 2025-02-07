package storage

import (
	"log"
	"os"
	"testing"

	"ethparser/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	logger := log.New(os.Stdout, "test: ", log.LstdFlags)
	storage := NewMemoryStorage(logger)

	t.Run("Subscribe", func(t *testing.T) {
		address := "0x123"

		// First subscription
		assert.True(t, storage.Subscribe(address))

		// Duplicate subscription
		assert.False(t, storage.Subscribe(address))

		// Verify subscription
		assert.True(t, storage.IsSubscribed(address))
	})

	t.Run("AddAndGetTransactions", func(t *testing.T) {
		address := "0x123"
		storage.Subscribe(address)

		tx := types.ParsedTransaction{
			Hash:        "0xabc",
			From:        address,
			To:          "0x456",
			Value:       "1000000000000000000",
			BlockNumber: 1000,
			Timestamp:   1612345678,
		}

		// Add transaction
		storage.AddTransaction(tx)

		// Get transactions
		txs := storage.GetTransactions(address)
		assert.Len(t, txs, 1)
		assert.Equal(t, tx.Hash, txs[0].Hash)
	})

	t.Run("CurrentBlock", func(t *testing.T) {
		// Set block
		blockNum := 1000
		storage.SetCurrentBlock(blockNum)

		// Get block
		assert.Equal(t, blockNum, storage.GetCurrentBlock())
	})
}
