// blockchain.go

package blockchain

import (
	"crypto/sha256"
	"time"

	"github.com/nnkienn/lab1-blc/merkletree"
)

type Transaction struct {
	Data       []byte
	MerkleRoot []byte
}

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	MerkleRoot    []byte
}

type Blockchain struct {
	Blocks []*Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{Blocks: []*Block{genesisBlock()}}
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := generateBlock(prevBlock, transactions)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) GetLatestBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *Blockchain) GetBlocks() []*Block {
	return bc.Blocks
}

func generateBlock(prevBlock *Block, transactions []*Transaction) *Block {
	newBlock := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlock.Hash,
	}

	// Calculate the Merkle root for transactions
	newBlock.MerkleRoot = buildMerkleTree(transactions).Root.Hash

	newBlock.SetHash()
	return newBlock
}

func (b *Block) SetHash() {
	headers := append(b.PrevBlockHash, b.HashTransaction()...)
	headers = append(headers, []byte(string(b.Timestamp))...)
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func (b *Block) HashTransaction() []byte {
	var transactionsData []byte
	for _, transaction := range b.Transactions {
		transactionsData = append(transactionsData, transaction.MerkleRoot...)
	}
	transactionHash := sha256.Sum256(transactionsData)
	return transactionHash[:]
}

func genesisBlock() *Block {
	return &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  []*Transaction{},
		PrevBlockHash: []byte{},
		MerkleRoot:    []byte{}, // Add a placeholder Merkle root for the genesis block
	}
}

func buildMerkleTree(transactions []*Transaction) *merkletree.MerkleTree {
	var txHashes [][]byte
	for _, tx := range transactions {
		txHashes = append(txHashes, tx.Data)
	}

	return merkletree.NewMerkleTree(txHashes)
}
