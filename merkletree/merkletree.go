package merkletree

import (
	"crypto/sha256"
)

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
