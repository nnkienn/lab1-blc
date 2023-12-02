package merkletree

import (
	"crypto/sha256"
)

// Transaction represents a basic transaction in the blockchain
type Transaction struct {
	Data []byte
}

// MerkleNode represents a node in the Merkle Tree
type MerkleNode struct {
	Hash        []byte
	Left, Right *MerkleNode
	Transaction *Transaction
}

// MerkleTree represents the Merkle Tree
type MerkleTree struct {
	Root *MerkleNode
}

// NewMerkleNode creates a new MerkleNode
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

// NewMerkleTree creates a new MerkleTree
func NewMerkleTree(transactions []*Transaction) *MerkleTree {
	var nodes []*MerkleNode

	for _, transaction := range transactions {
		hash := calculateHash(transaction.Data)
		nodes = append(nodes, &MerkleNode{Hash: hash, Transaction: transaction})
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

// calculateHash calculates the SHA-256 hash of the data
func calculateHash(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// GetMerkleTreeData returns the data of the MerkleTree
func (mt *MerkleTree) GetMerkleTreeData() []*Transaction {
	var transactions []*Transaction
	mt.traverse(mt.Root, &transactions)
	return transactions
}

// traverse traverses the MerkleTree and collects transaction data
func (mt *MerkleTree) traverse(node *MerkleNode, transactions *[]*Transaction) {
	if node == nil {
		return
	}

	if node.Left == nil && node.Right == nil {
		*transactions = append(*transactions, node.Transaction)
		return
	}

	mt.traverse(node.Left, transactions)
	mt.traverse(node.Right, transactions)
}
