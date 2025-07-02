package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob" // Import encoding/gob
	"encoding/hex"
	"log"
	"os" // Import os

	"github.com/btcsuite/btcd/btcutil/base58" // Import base58 library
	"golang.org/x/crypto/ripemd160"           // For generating addresses
)

const walletFile = "wallets.dat" // Define wallet file name

// Wallet represents a cryptocurrency wallet
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

// Wallets represents a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// NewWallets creates and returns a new Wallets instance
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile() // Call the method on the wallets instance
	if err != nil && !os.IsNotExist(err) {
		// If the file doesn't exist, it's not an error, just means no wallets saved yet.
		// If it's another error, return it.
		log.Printf("Error loading wallets from file: %v", err)
		return nil, err
	}

	return &wallets, nil // Return wallets (either loaded or empty) and nil error
}

// NewWallet creates and returns a new Wallet
func NewWallet() *Wallet {
	private, public := newKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

// newKeyPair creates a new public and private key pair
func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return private, pubKey
}

// GetAddress returns the wallet address
func (w Wallet) GetAddress() []byte {
	pubHash := PublicKeyHash(w.PublicKey)

	version := byte(0x00) // Mainnet version byte
	payload := append([]byte{version}, pubHash...)
	checksum := checksum(payload)

	fullPayload := append(payload, checksum...)
	address := base58.Encode(fullPayload) // Use base58.Encode

	return []byte(address) // Return as byte slice
}

// PublicKeyHash hashes the public key
func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := hasher.Sum(nil)

	return publicRIPEMD160
}

// checksum generates a checksum for a public key hash
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:4]
}

// Base58Decode decodes a Base58 encoded byte slice
func Base58Decode(input []byte) []byte {
	return base58.Decode(string(input))
}

// ValidateAddress checks if address is valid
func ValidateAddress(address []byte) bool {
	pubKeyHash := Base58Decode(address)
	if len(pubKeyHash) < 4 { // Check minimum length
		return false
	}

	checksumBytes := pubKeyHash[len(pubKeyHash)-4:]
	versionedPayload := pubKeyHash[:len(pubKeyHash)-4]
	actualChecksum := checksum(versionedPayload)

	return hex.EncodeToString(checksumBytes) == hex.EncodeToString(actualChecksum)
}

// SaveToFile saves the wallets to a file
func (ws *Wallets) SaveToFile() {
	file, err := os.Create(walletFile)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}
}

// LoadFromFile loads wallets from a file
func (ws *Wallets) LoadFromFile() error {
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err // File does not exist, return the error
	}

	file, err := os.Open(walletFile)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(ws)
	if err != nil {
		log.Panic(err)
	}

	return nil // Return nil on success
}