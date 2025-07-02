package api

import (
	"fmt" // Import fmt package
	"net/http"
	"time" // Import time package

	"github.com/gin-gonic/gin"

	"aztecs/consensus" // Import consensus package
	"aztecs/core" // Import core package
	"aztecs/crypto" // Import crypto package
)

// RegisterRoutes registers the API routes
func RegisterRoutes(router *gin.Engine, bc *core.Blockchain, wallets *crypto.Wallets) { // Accept Blockchain and Wallets instances
	router.GET("/blockchain", func(c *gin.Context) {
		getBlockchain(c, bc) // Pass context and blockchain instance
	})
	router.POST("/mine", func(c *gin.Context) {
		mineBlock(c, bc) // Pass context and blockchain instance
	})
	router.POST("/transactions", func(c *gin.Context) { // Use anonymous function
		createTransaction(c, bc) // Pass context and blockchain instance
	})
	router.POST("/wallets", func(c *gin.Context) { // Use anonymous function
		createWallet(c, wallets) // Pass context and wallets instance
	})
	router.GET("/wallets", func(c *gin.Context) { // Add get wallets route
		getWallets(c, wallets) // Pass context and wallets instance
	})
	router.GET("/wallets/:address", getWallet) // TODO: Pass wallets instance
	router.GET("/wallets/:address/balance", func(c *gin.Context) { // Add get wallet balance route
		getWalletBalance(c, bc) // Pass context and blockchain instance
	})
}

// getWallets handles the request to get all wallets
func getWallets(c *gin.Context, wallets *crypto.Wallets) { // Accept Wallets instance
	addresses := []string{}
	for address := range wallets.Wallets {
		addresses = append(addresses, address)
	}
	c.JSON(http.StatusOK, gin.H{"wallets": addresses})
}

// getBlockchain handles the request to get the blockchain
func getBlockchain(c *gin.Context, bc *core.Blockchain) { // Accept Blockchain instance
	c.JSON(http.StatusOK, bc.Blocks) // Return the blockchain blocks
}

// mineBlock handles the request to mine a new block
func mineBlock(c *gin.Context, bc *core.Blockchain) { // Accept Blockchain instance
	// Create a new block with placeholder data for now
	newBlock := core.NewBlock(int64(len(bc.Blocks)), time.Now(), "Block mined via API", bc.Blocks[len(bc.Blocks)-1].Hash) // Need to import time

	// Perform Proof of Work
	pow := consensus.NewProofOfWork(newBlock) // Need to import consensus
	nonce, hash := pow.Run()

	newBlock.Hash = hash
	newBlock.Nonce = nonce

	// Add the mined block to the blockchain
	bc.Blocks = append(bc.Blocks, newBlock)

	c.JSON(http.StatusOK, gin.H{"message": "Block mined successfully", "block": newBlock})
}

// createTransaction handles the request to create a new transaction
func createTransaction(c *gin.Context, bc *core.Blockchain) { // Accept Blockchain instance
	var tx core.Transaction // Assuming Transaction struct is in core package
	if err := c.ShouldBindJSON(&tx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual transaction validation (e.g., sufficient balance, valid signature)

	// For simplicity, add the transaction data to the next block's data
	// In a real scenario, transactions would go into a transaction pool
	// and be included in a block during mining.
	// Here, we'll just append the transaction details to the last block's data for demonstration.
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	lastBlock.Data += fmt.Sprintf("\nTransaction: From %s To %s Amount %.2f", tx.From, tx.To, tx.Amount)

	c.JSON(http.StatusOK, gin.H{"message": "Transaction received and added to last block (placeholder)", "transaction": tx})
}

// createWallet handles the request to create a new wallet
func createWallet(c *gin.Context, wallets *crypto.Wallets) { // Accept Wallets instance
	wallet := crypto.NewWallet() // Create a new wallet
	address := wallet.GetAddress() // Get wallet address

	wallets.Wallets[string(address)] = wallet // Add wallet to manager

	// TODO: Save wallets to file

	c.JSON(http.StatusOK, gin.H{"address": string(address)})
}

// getWallet handles the request to get wallet details
// TODO: Implement actual logic to get wallet details
func getWallet(c *gin.Context) {
	address := c.Param("address")
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Get wallet endpoint for address %s (placeholder)", address)})
}

// getWalletBalance handles the request to get wallet balance
func getWalletBalance(c *gin.Context, bc *core.Blockchain) { // Accept Blockchain instance
	address := c.Param("address")
	balance := bc.GetBalance(address) // Call GetBalance on the blockchain instance
	c.JSON(http.StatusOK, gin.H{"address": address, "balance": balance})
}