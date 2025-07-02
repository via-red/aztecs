package consensus

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"

	"aztecs/core"
)

const targetBits = 16 // Difficulty level (adjust as needed)

// ProofOfWork represents a Proof of Work system
type ProofOfWork struct {
	block  *core.Block
	target *big.Int
}

// NewProofOfWork creates a new ProofOfWork instance
func NewProofOfWork(b *core.Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	pow := &ProofOfWork{b, target}
	return pow
}

// prepareData prepares data for hashing
func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			IntToHex(pow.block.Index),
			IntToHex(pow.block.Timestamp.Unix()), // Convert time to Unix timestamp
			[]byte(pow.block.Data),
			[]byte(pow.block.PrevHash),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

// Run performs the Proof of Work mining
func (pow *ProofOfWork) Run() (int, string) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining a new block with data: \"%s\"\n", pow.block.Data)
	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data)
		// fmt.Printf("\r%x", hash) // Optional: print hash during mining
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			fmt.Printf("Block mined! Hash: %x\n", hash)
			break
		} else {
			nonce++
		}
	}
	return nonce, hex.EncodeToString(hash[:])
}

// Validate validates the block's Proof of Work
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.target) == -1
}

// IntToHex converts an int64 to a byte slice
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}