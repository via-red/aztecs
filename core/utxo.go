package core

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

const utxoFile = "utxo.dat" // Define UTXO data file name

// UTXO represents an unspent transaction output
type UTXO struct {
	TxID      string  // ID of the transaction the output belongs to
	Index     int     // Index of the output in the transaction
	Value     float64 // Value of the output
	PubKeyHash []byte  // Public key hash of the recipient (raw bytes)
}

// UTXOSet represents the collection of unspent transaction outputs
// Map: TxID -> Map: Vout -> UTXO
type UTXOSet struct {
	UTXOs map[string]map[int]*UTXO
}

// NewUTXOSet creates a new UTXOSet or loads from file
func NewUTXOSet() (*UTXOSet, error) {
	utxoSet := UTXOSet{}
	// Initialize the outer map
	utxoSet.UTXOs = make(map[string]map[int]*UTXO)

	err := utxoSet.LoadFromFile()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading UTXO set from file: %v", err)
		return nil, err
	}

	// If loaded UTXOs map is nil (e.g., file was empty or corrupted), re-initialize
	if utxoSet.UTXOs == nil {
		utxoSet.UTXOs = make(map[string]map[int]*UTXO)
	}

	// Note: The initial UTXO set is built from the blockchain in core/chain.go
	// if the utxoFile does not exist or is empty after loading.

	return &utxoSet, nil
}

// SaveToFile saves the UTXO set to a file
func (uset *UTXOSet) SaveToFile() {
	file, err := os.Create(utxoFile)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(uset)
	if err != nil {
		log.Panic(err)
	}
	log.Println("UTXO set saved to file.")
}

// LoadFromFile loads the UTXO set from a file
func (uset *UTXOSet) LoadFromFile() error {
	if _, err := os.Stat(utxoFile); os.IsNotExist(err) {
		return err // File does not exist
	}

	file, err := os.Open(utxoFile)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(uset)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

// FindUTXOs finds all UTXOs for a given public key hash
func (uset *UTXOSet) FindUTXOs(pubKeyHash []byte) []*UTXO {
	foundUTXOs := []*UTXO{}

	// Iterate through all transactions in the UTXO set
	for _, vouts := range uset.UTXOs {
		// Iterate through all outputs of the transaction
		for _, utxo := range vouts {
			// Compare public key hashes directly
			if bytes.Equal(utxo.PubKeyHash, pubKeyHash) {
				foundUTXOs = append(foundUTXOs, utxo)
			}
		}
	}
	return foundUTXOs
}

// BuildFromBlockchain builds the UTXO set by scanning the entire blockchain
func (uset *UTXOSet) BuildFromBlockchain(bc *Blockchain) {
	// Clear the existing UTXO set and re-initialize the outer map
	uset.UTXOs = make(map[string]map[int]*UTXO)

	// Map to track spent transaction outputs (TxID -> list of output indices)
	spentTXOs := make(map[string][]int)

	// Iterate through all blocks in the blockchain
	// We iterate in reverse order to easily identify spent outputs
	for i := len(bc.Blocks) - 1; i >= 0; i-- {
		block := bc.Blocks[i]

		// Iterate through transactions in the block
		for _, tx := range block.Transactions {
			// Process transaction outputs
			// Add unspent outputs to the UTXO set
			if !tx.IsCoinbase() {
				// For regular transactions, mark inputs as spent
				for _, vin := range tx.Vin {
					txid := vin.Txid
					spentTXOs[txid] = append(spentTXOs[txid], vin.Vout)
				}
			}

			// Process transaction outputs
			// Add outputs to the UTXO set if they are not spent
			txid := tx.ID
			for voutIndex, vout := range tx.Vout {
				// Check if the output has been spent
				if spentTXOs[txid] != nil {
					spent := false
					for _, spentVout := range spentTXOs[txid] {
						if spentVout == voutIndex {
							spent = true
							break
						}
					}
					if spent {
						continue // Output is spent, skip
					}
				}

				// Output is unspent, add it to the UTXO set
				utxo := UTXO{
					TxID:      txid,
					Index:     voutIndex,
					Value:     vout.Value,
					PubKeyHash: vout.PubKeyHash, // Store raw public key hash
				}

				// Initialize the inner map if it doesn't exist for this TxID
				if uset.UTXOs[txid] == nil {
					uset.UTXOs[txid] = make(map[int]*UTXO)
				}
				uset.UTXOs[txid][voutIndex] = &utxo // Add the UTXO to the set
			}
		}
	}
	log.Println("UTXO set built from blockchain.")
}

// Update updates the UTXO set based on new blocks
func (uset *UTXOSet) Update(block *Block) {
	// Process transaction inputs (mark spent outputs)
	for _, tx := range block.Transactions {
		if !tx.IsCoinbase() {
			for _, vin := range tx.Vin {
				// Remove the spent UTXO from the set
				// Check if the transaction ID exists in the UTXO set
				if uset.UTXOs[vin.Txid] != nil {
					// Check if the output index exists for this transaction ID
					if uset.UTXOs[vin.Txid][vin.Vout] != nil {
						// Remove the UTXO
						delete(uset.UTXOs[vin.Txid], vin.Vout)
						log.Printf("Removed spent UTXO %s:%d from UTXO set", vin.Txid, vin.Vout)

						// If there are no more outputs for this transaction ID, remove the inner map
						if len(uset.UTXOs[vin.Txid]) == 0 {
							delete(uset.UTXOs, vin.Txid)
							log.Printf("Removed empty inner map for TxID %s", vin.Txid)
						}
					} else {
						log.Printf("Warning: Attempted to remove non-existent UTXO %s:%d", vin.Txid, vin.Vout)
					}
				} else {
					log.Printf("Warning: Attempted to remove UTXO from non-existent TxID %s", vin.Txid)
				}
			}
		}

		// Process transaction outputs (add new UTXOs)
		txid := tx.ID
		for voutIndex, vout := range tx.Vout {
			utxo := UTXO{
				TxID:      txid,
				Index:     voutIndex,
				Value:     vout.Value,
				PubKeyHash: vout.PubKeyHash, // Store raw public key hash
			}

			// Initialize the inner map if it doesn't exist for this TxID
			if uset.UTXOs[txid] == nil {
				uset.UTXOs[txid] = make(map[int]*UTXO)
			}
			uset.UTXOs[txid][voutIndex] = &utxo // Add the new UTXO
			log.Printf("Added new UTXO %s:%d to UTXO set", txid, voutIndex)
		}
	}
	log.Println("UTXO set updated.")
}