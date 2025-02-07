package types

// Block represents an Ethereum block structure
type Block struct {
	Number       string        `json:"number"`
	Hash         string        `json:"hash"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
}

// Transaction represents the raw transaction from Ethereum RPC
type Transaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber string `json:"blockNumber"`
}

// ParsedTransaction represents our processed transaction with converted values
type ParsedTransaction struct {
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	BlockNumber int64  `json:"blockNumber"`
	Timestamp   int64  `json:"timestamp"`
}

type Parser interface {
	// GetCurrentBlock - last parsed block
	GetCurrentBlock() int

	// Subscribe - add address to observer
	Subscribe(address string) bool

	// GetTransactions - list of inbound or outbound transactions for an address
	GetTransactions(address string) []ParsedTransaction
}
