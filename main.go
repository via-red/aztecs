package main

import (
	"aztecs/api"
	"aztecs/core"
	"aztecs/storage"
	"fmt"
)

func main() {
	fmt.Println("Simple Blockchain Project")

	// Initialize database
	db := storage.NewBlockchainDB()
	defer db.Close()

	// Initialize blockchain
	bc, err := core.NewBlockchain()
	if err != nil {
		fmt.Println("Error initializing blockchain:", err)
		return // Exit if blockchain initialization fails
	}

	// Initialize wallet manager
	// Wallets are now initialized within api.StartServer, which calls crypto.NewWallets
	// wallets, err := crypto.NewWallets()
	// if err != nil {
	// 	fmt.Println("Error initializing wallets:", err)
	// 	return // Exit if wallets initialization fails
	// }

	// Add a few blocks with placeholder transactions for testing
	// In a real application, transactions would be created and added via API
	if len(bc.Blocks) == 1 { // Only add blocks if it's a new blockchain with just the genesis block
		// Create placeholder transactions using the new structure
		tx1 := &core.Transaction{
			ID: "tx1", // Placeholder ID
			Vin: []core.TxInput{
				{Txid: "prev_tx_id_A", Vout: 0, Signature: nil, PubKey: []byte("addressA_pubkey")}, // Placeholder input
			},
			Vout: []core.TxOutput{
				{Value: 10.0, PubKeyHash: []byte("addressB_hash")}, // Placeholder output
			},
		}
		tx1.SetID() // Set the transaction ID

		tx2 := &core.Transaction{
			ID: "tx2", // Placeholder ID
			Vin: []core.TxInput{
				{Txid: "prev_tx_id_B", Vout: 0, Signature: nil, PubKey: []byte("addressB_pubkey")}, // Placeholder input
			},
			Vout: []core.TxOutput{
				{Value: 5.0, PubKeyHash: []byte("addressC_hash")}, // Placeholder output
			},
		}
		tx2.SetID() // Set the transaction ID

		bc.AddBlock([]*core.Transaction{tx1}) // Add block with transaction
		bc.AddBlock([]*core.Transaction{tx2}) // Add block with transaction
	}


	// Check blockchain validity (optional, can be removed later)
	// fmt.Printf("Blockchain is valid: %v\n", bc.IsValid())

	// Start API server, passing the blockchain instance
	api.StartServer(bc) // Pass the blockchain instance
}