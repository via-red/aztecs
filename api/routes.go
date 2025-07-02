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
	// Create a new block with empty transactions
	newBlock := core.NewBlock(int64(len(bc.Blocks)), time.Now(), []*core.Transaction{}, bc.Blocks[len(bc.Blocks)-1].Hash)

	// Perform Proof of Work
	pow := consensus.NewProofOfWork(newBlock)
	nonce, hash := pow.Run()

	newBlock.Hash = hash
	newBlock.Nonce = nonce

	// Add the mined block to the blockchain
	bc.Blocks = append(bc.Blocks, newBlock)

	c.JSON(http.StatusOK, gin.H{"message": "Block mined successfully", "block": newBlock})
}

// createTransaction handles the request to create a new transaction
func createTransaction(c *gin.Context, bc *core.Blockchain) { // Accept Blockchain instance
	// 使用临时结构体接收前端数据
	type TempTx struct {
		From   string  `json:"fromAddress"`
		To     string  `json:"toAddress"`
		Amount float64 `json:"amount"`
	}
	
	var tempTx TempTx
	if err := c.ShouldBindJSON(&tempTx); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建实际交易对象（简化版）
	tx := &core.Transaction{
		Vin: []core.TxInput{
			{Txid: "prev_tx", Vout: 0, PubKey: []byte(tempTx.From)},
		},
		Vout: []core.TxOutput{
			{Value: tempTx.Amount, PubKeyHash: []byte(tempTx.To)},
		},
	}
	tx.SetID()

	// TODO: 实际实现应添加交易到交易池
	c.JSON(http.StatusOK, gin.H{"message": "Transaction received (placeholder)", "transaction": tx})
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
	balance := bc.GetBalance([]byte(address)) // Convert address to []byte
	c.JSON(http.StatusOK, gin.H{"address": address, "balance": balance})
}