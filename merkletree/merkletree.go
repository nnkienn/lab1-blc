package merkletree

import (
    "crypto/sha256"
)

// Transaction đại diện cho một giao dịch cơ bản trong blockchain
type Transaction struct {
    Data []byte
}

// MerkleNode đại diện cho một nút trong cây Merkle
type MerkleNode struct {
    Hash        []byte
    Left, Right *MerkleNode
    Transaction *Transaction
}

// MerkleTree đại diện cho cây Merkle
type MerkleTree struct {
    Root *MerkleNode
}

// NewMerkleNode tạo một nút MerkleNode mới
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

// NewMerkleTree tạo một cây Merkle mới
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

// calculateHash tính toán hash SHA-256 của dữ liệu
func calculateHash(data []byte) []byte {
    hash := sha256.New()
    hash.Write(data)
    return hash.Sum(nil)
}
