package storage

import (
	"fmt" // Import fmt package
	"log"

	"github.com/boltdb/bolt" // Using BoltDB for storage
)

const dbFile = "blockchain.db"
const blocksBucket = "BlocksBucket"

// BlockchainDB represents the database connection
type BlockchainDB struct {
	db *bolt.DB
}

// NewBlockchainDB creates or opens the BoltDB database
func NewBlockchainDB() *BlockchainDB {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// Create the blocks bucket if it doesn't exist
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(blocksBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err) // Need to import fmt
		}
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &BlockchainDB{db: db}
}

// SaveBlock saves a block to the database
// TODO: Implement saving block data
func (bdb *BlockchainDB) SaveBlock(block []byte) {
	// Placeholder for saving block data
	log.Println("Saving block to DB (placeholder)")
}

// GetLastBlock gets the last block hash from the database
// TODO: Implement getting last block hash
func (bdb *BlockchainDB) GetLastBlock() []byte {
	// Placeholder for getting last block hash
	log.Println("Getting last block from DB (placeholder)")
	return nil // Return nil for now
}

// Close closes the database connection
func (bdb *BlockchainDB) Close() {
	err := bdb.db.Close()
	if err != nil {
		log.Panic(err)
	}
}