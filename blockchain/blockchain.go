package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// Transaction represents a single transaction in the blockchain
type Transaction struct {
	Data []byte
}

// MerkleNode represents a node in the Merkle tree
type MerkleNode struct {
	Hash  []byte
	Left  *MerkleNode
	Right *MerkleNode
}

// MerkleTree represents the Merkle tree
type MerkleTree struct {
	Root *MerkleNode
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	var hash []byte
	if left == nil && right == nil {
		hash = calculateHash(data)
	} else {
		hashData := append(left.Hash, right.Hash...)
		hash = calculateHash(hashData)
	}

	return &MerkleNode{Hash: hash, Left: left, Right: right}
}

// NewMerkleTree creates a new Merkle tree from a list of transaction hashes
func NewMerkleTree(transactionHashes [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	for _, hash := range transactionHashes {
		nodes = append(nodes, &MerkleNode{Hash: hash})
	}

	for len(nodes) > 1 {
		var level []*MerkleNode
		for i := 0; i < len(nodes); i += 2 {
			var left, right *MerkleNode
			if i+1 < len(nodes) {
				left = nodes[i]
				right = nodes[i+1]
			} else {
				left = nodes[i]
				right = nil
			}
			node := NewMerkleNode(left, right, nil)
			level = append(level, node)
		}
		nodes = level
	}

	return &MerkleTree{Root: nodes[0]}
}

// calculateHash calculates the SHA-256 hash of the input data
func calculateHash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
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
