package core

import (
	"encoding/gob" // Import encoding/gob
	"fmt"          // Import fmt
	"log"
	"os" // Import os
	"time"
)

const blockchainFile = "blockchain.dat" // Define blockchain data file name

// Blockchain represents the blockchain
type Blockchain struct {
	Blocks  []*Block
	UTXOSet *UTXOSet // Add UTXO set to the blockchain
}

// NewBlockchain creates a new blockchain with a genesis block or loads from file
func NewBlockchain() (*Blockchain, error) {
	// Check if blockchain data file exists
	if _, err := os.Stat(blockchainFile); os.IsNotExist(err) {
		// File does not exist, create a new blockchain with genesis block
		log.Println("Blockchain data file not found, creating new blockchain.")
		// Create a placeholder genesis transaction
		// Create a placeholder genesis transaction (coinbase transaction)
		// Coinbase transaction has no inputs and one output
		genesisTx := &Transaction{
			ID: "genesis_tx", // Placeholder ID
			Vin: []TxInput{
				{Txid: "", Vout: -1, Signature: nil, PubKey: []byte("coinbase")}, // Coinbase input
			},
			Vout: []TxOutput{
				{Value: 50.0, PubKeyHash: []byte("genesis_address_hash")}, // Output to genesis address (placeholder hash)
			},
		}
		// Set the transaction ID
		genesisTx.SetID() // Calculate and set the transaction ID

		genesisBlock := NewBlock(0, time.Now(), []*Transaction{genesisTx}, "") // Pass a slice with the genesis transaction
		// For simplicity, we'll calculate the genesis block hash directly here
		// In a real scenario, mining would be involved
		genesisBlock.Hash = genesisBlock.CalculateHash()
		bc := &Blockchain{Blocks: []*Block{genesisBlock}}

		// Create and build the initial UTXO set from the genesis block
		utxoSet, err := NewUTXOSet() // Create a new empty UTXO set and handle error
		if err != nil {
			return nil, fmt.Errorf("failed to create new UTXO set: %w", err)
		}
		// TODO: Implement BuildFromBlockchain method in utxo.go
		// utxoSet.BuildFromBlockchain(bc) // Build from the blockchain
		bc.UTXOSet = utxoSet // Assign the UTXO set to the blockchain
		utxoSet.SaveToFile() // Save the initial UTXO set

		bc.SaveToFile() // Save the newly created blockchain
		return bc, nil
	}

	// File exists, load blockchain from file
	log.Println("Blockchain data file found, loading blockchain.")
	file, err := os.Open(blockchainFile) // Use = for assignment
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	var bc Blockchain
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&bc) // Use = for assignment
	if err != nil {
		log.Panic(err)
	}

	// Load or build the UTXO set
	utxoSet, err := NewUTXOSet() // Try to load UTXO set and handle error
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading UTXO set: %v", err)
		return nil, err // Return error if loading fails for reasons other than file not existing
	}

	if utxoSet.UTXOs == nil || len(utxoSet.UTXOs) == 0 {
		// If UTXO set file didn't exist or was empty, build it from the loaded blockchain
		log.Println("UTXO set data file not found or empty, building from blockchain.")
		// TODO: Implement BuildFromBlockchain method in utxo.go
		// utxoSet.BuildFromBlockchain(&bc) // Build from the loaded blockchain
		utxoSet.SaveToFile() // Save the newly built UTXO set
	} else {
		log.Println("UTXO set loaded successfully.")
	}

	bc.UTXOSet = utxoSet // Assign the UTXO set to the blockchain

	log.Println("Blockchain loaded successfully.")
	return &bc, nil // Return loaded blockchain and nil error
}

// AddBlock adds a new block to the blockchain
// TODO: Accept a slice of Transactions instead of a data string
func (bc *Blockchain) AddBlock(transactions []*Transaction) { // Accept a slice of Transactions
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Index+1, time.Now(), transactions, prevBlock.Hash) // Pass the transactions
	// In a real scenario, mining would happen here to find a valid hash
	newBlock.Hash = newBlock.CalculateHash() // Placeholder for mining
	bc.Blocks = append(bc.Blocks, newBlock)
	log.Printf("Block #%d added to the blockchain", newBlock.Index)
	// TODO: Update UTXO set based on the actual transactions in the block
	// bc.UTXOSet.Update(newBlock) // Update the UTXO set with the new block
	// bc.UTXOSet.SaveToFile() // Save the updated UTXO set
	bc.SaveToFile() // Save the blockchain after adding a block
}

// IsValid checks if the blockchain is valid
func (bc *Blockchain) IsValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		// Check if the current block's hash is correct
		if currentBlock.Hash != currentBlock.CalculateHash() {
			log.Println("Block hash mismatch")
			return false
		}

		// Check if the current block's previous hash matches the previous block's hash
		if currentBlock.PrevHash != prevBlock.Hash {
			log.Println("Previous block hash mismatch")
			return false
		}
	}
	return true
}

// GetBalance calculates the balance of a given address
// TODO: Implement actual balance calculation using UTXO set
// GetBalance calculates the balance for a given public key hash
func (bc *Blockchain) GetBalance(pubKeyHash []byte) float64 {
	balance := 0.0
	utxos := bc.UTXOSet.FindUTXOs(pubKeyHash)
	for _, utxo := range utxos {
		balance += utxo.Value
	}
	log.Printf("Balance for pubKeyHash: %x: %.2f", pubKeyHash, balance)
	return balance
}

// SaveToFile saves the blockchain to a file
func (bc *Blockchain) SaveToFile() {
	file, err := os.Create(blockchainFile)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(bc)
	if err != nil {
		log.Panic(err)
	}
	log.Println("Blockchain saved to file.")
}

// LoadFromFile loads the blockchain from a file
// This method is now integrated into NewBlockchain, but kept for potential future use or clarity.
/*
func (bc *Blockchain) LoadFromFile() error {
	if _, err := os.Stat(blockchainFile); os.IsNotExist(err) {
		return err // File does not exist
	}

	file, err := os.Open(blockchainFile)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(bc)
	if err != nil {
		log.Panic(err)
	}

	return nil
}
*/
