package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
	"github.com/nnkienn/lab1-blc/merkletree"
)

// Transaction represents a single transaction in the blockchain
type Transaction struct {
	Data []byte
}

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

// NewTransaction creates a new transaction
func NewTransaction(data []byte) *Transaction {
	return &Transaction{Data: data}
}

// NewBlock creates a new block in the blockchain
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}

	// Calculate Merkle Root
	block.calculateMerkleRoot()

	// Calculate Block Hash
	block.calculateHash()

	return block
}

// calculateMerkleRoot calculates the Merkle Root of a block's transactions
func (block *Block) calculateMerkleRoot() {
	var transactionData [][]byte

	for _, tx := range block.Transactions {
		transactionData = append(transactionData, tx.Data)
	}

	merkleTree := NewMerkleTree(transactionData)
	block.MerkleRoot = merkleTree.Root.Hash
}

// calculateHash calculates the hash of the block
func (block *Block) calculateHash() {
	header := append(block.PrevBlockHash, block.MerkleRoot...)
	header = append(header, []byte(string(block.Timestamp))...)
	hash := sha256.Sum256(header)
	block.Hash = hash[:]
}

// HexHash returns the hexadecimal representation of the block's hash
func (block *Block) HexHash() string {
	return hex.EncodeToString(block.Hash)
}
