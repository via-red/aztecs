package api

import (
	"log"

	"github.com/gin-gonic/gin" // Using Gin framework

	"aztecs/core" // Import core package
	"aztecs/crypto" // Import crypto package
)

// StartServer starts the HTTP API server
func StartServer(bc *core.Blockchain) { // Accept Blockchain instance
	router := gin.Default()

	// Initialize wallet manager with error handling
	wallets, err := crypto.NewWallets()
	if err != nil {
		log.Fatalf("Failed to initialize wallets: %v", err)
	}

	// Define API routes
	RegisterRoutes(router, bc, wallets) // Pass Blockchain and Wallets instances to routes

	log.Println("Starting API server on :8080")
	err = router.Run(":8080")
	if err != nil {
		log.Panic(err)
	}
}