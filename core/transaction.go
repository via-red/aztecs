package core

import (
	"bytes" // Import bytes
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/gob" // Import encoding/gob
	"encoding/hex"
	"log"

	"aztecs/crypto" // Import crypto package
)

// TxInput represents a transaction input
type TxInput struct {
	Txid      string // ID of the transaction the output is from
	Vout      int    // Index of the output in the transaction
	Signature []byte // Signature to unlock the output
	PubKey    []byte // Public key of the sender
}

// TxOutput represents a transaction output
type TxOutput struct {
	Value   float64 // Value of the output
	PubKeyHash []byte // Hash of the recipient's public key
}

// Transaction represents a transaction in the blockchain
type Transaction struct {
	ID   string
	Vin  []TxInput  // Transaction inputs
	Vout []TxOutput // Transaction outputs
}

// CalculateHash calculates the hash of the transaction
// SetID calculates and sets the ID of the transaction
func (tx *Transaction) SetID() {
	// TODO: Include Vin and Vout in the hash calculation
	// For now, a simplified hash
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hex.EncodeToString(hash[:])
}

// CalculateHash calculates the hash of the transaction
// This method is similar to SetID but returns the hash instead of setting the ID.
// It might be redundant if SetID is always called after creating a transaction.
// Let's keep it for now but consider if it's needed.
func (tx *Transaction) CalculateHash() string {
	// TODO: Include Vin and Vout in the hash calculation
	// For now, a simplified hash
	var encoded bytes.Buffer
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash := sha256.Sum256(encoded.Bytes())
	return hex.EncodeToString(hash[:])
}

// Sign signs the transaction with the private key
// TODO: Implement actual signing with ECC
// Sign signs the transaction with the private key
// TODO: Implement actual signing of transaction inputs
func (tx *Transaction) Sign(privateKey *ecdsa.PrivateKey) {
	// This signing logic needs to be updated to sign transaction inputs.
	// For now, this is a placeholder.
	log.Println("Transaction signing is a placeholder.")
	// Example (illustrative, needs proper implementation):
	/*
		txCopy := tx.TrimmedCopy() // Create a copy without signatures and public keys
		for i, vin := range txCopy.Vin {
			// Get the public key hash of the output being spent by this input
			// This requires looking up the previous transaction (vin.Txid)
			// and the specific output (vin.Vout).
			// This logic depends on having access to the blockchain or UTXO set.

			// For now, let's use a placeholder signature
			signature := []byte("placeholder_signature") // Replace with actual signature
			tx.Vin[i].Signature = signature
			tx.Vin[i].PubKey = privateKey.PublicKey.X.Bytes() // Placeholder for public key
		}
	*/
}

// IsValid checks if the transaction is valid
// TODO: Implement actual verification with ECC
// IsValid checks if the transaction is valid
// TODO: Implement actual transaction validation using UTXO set
func (tx *Transaction) IsValid() bool {
	// This validation logic needs to be updated to verify transaction inputs.
	// For now, this is a placeholder.
	log.Println("Transaction validation is a placeholder.")
	// In a real implementation, you would:
	// - Verify that the sum of input values equals the sum of output values (minus fees).
	// - Verify the signature of each input against the public key hash in the corresponding output.
	// - Ensure that the inputs are actually unspent outputs in the UTXO set.

	// Placeholder: Always return true for now
	return true
}

// IsCoinbase checks if a transaction is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	// A coinbase transaction has only one input, and its Txid is empty
	// and Vout is -1.
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// TrimmedCopy creates a trimmed copy of the transaction for signing
// This copy excludes signatures and public keys from inputs.
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TxInput
	for _, vin := range tx.Vin {
		inputs = append(inputs, TxInput{Txid: vin.Txid, Vout: vin.Vout, Signature: nil, PubKey: nil})
	}
	var outputs []TxOutput
	for _, vout := range tx.Vout {
		outputs = append(outputs, TxOutput{Value: vout.Value, PubKeyHash: vout.PubKeyHash})
	}
	return Transaction{ID: tx.ID, Vin: inputs, Vout: outputs}
}

// UsesKey checks if the input's public key hash is the same as the provided public key hash
// UsesKey checks if the input's public key hash is the same as the provided public key hash
// UsesKey checks if the input's public key hash is the same as the provided public key hash
func (in *TxInput) UsesKey(pubKeyHash []byte) bool {
	// Hash the input's PubKey and compare it to the provided pubKeyHash
	inputPubKeyHash := crypto.PublicKeyHash(in.PubKey) // Use the actual PublicKeyHash function

	return bytes.Compare(inputPubKeyHash, pubKeyHash) == 0
}

// IsLockedWithKey checks if the output's public key hash is the same as the provided public key hash
// IsLockedWithKey checks if the output's public key hash is the same as the provided public key hash
func (out *TxOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	// Compare the output's PubKeyHash with the provided pubKeyHash
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}