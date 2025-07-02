package core

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Index     int64
	Timestamp time.Time
	Transactions []*Transaction // Replace Data with a slice of Transactions
	PrevHash  string
	Hash      string
	Nonce     int
}

// CalculateHash calculates the hash of the block
// CalculateHash calculates the hash of the block
// TODO: Include transactions in the hash calculation
func (b *Block) CalculateHash() string {
	// For now, we'll exclude transactions from the hash calculation.
	// In a real blockchain, you would hash the transactions (e.g., using a Merkle tree).
	record := string(b.Index) + b.Timestamp.String() + b.PrevHash + string(b.Nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// NewBlock creates a new block
// TODO: Accept a slice of Transactions instead of a data string
func NewBlock(index int64, timestamp time.Time, transactions []*Transaction, prevHash string) *Block {
	block := &Block{
		Index:     index,
		Timestamp: timestamp,
		Transactions: transactions, // Assign the transactions
		PrevHash:  prevHash,
		Nonce:     0, // Initial nonce
	}
	// Hash will be calculated during mining (PoW)
	return block
}