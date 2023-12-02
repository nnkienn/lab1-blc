package blockchain

import (
	"fmt"
	"merkletree/merkletree"
	"p2p/p2p"
)

// Transaction represents a basic transaction in the blockchain
type Transaction struct {
	Data []byte
}

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
}

// Blockchain represents the entire blockchain
type Blockchain struct {
	blocks []*Block
}

// NewBlockchain creates a new blockchain with a genesis block
func NewBlockchain() *Blockchain {
	genesisTransaction := &Transaction{Data: []byte("Genesis Transaction")}
	genesisBlock := NewBlock([]byte{}, []*Transaction{genesisTransaction})
	return &Blockchain{blocks: []*Block{genesisBlock}}
}

// NewBlock creates a new block with the given transactions and previous block hash
func NewBlock(prevBlockHash []byte, transactions []*Transaction) *Block {
	block := &Block{
		Timestamp:     getCurrentTimestamp(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}
	block.SetHash()
	return block
}

// getCurrentTimestamp returns the current timestamp
func getCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// SetHash calculates and sets the hash of the block based on its fields
func (b *Block) SetHash() {
	headers := append(b.PrevBlockHash, getTransactionsHash(b.Transactions)...)
	hash := sha256.New()
	hash.Write(headers)
	b.Hash = hash.Sum(nil)
}

// getTransactionsHash calculates the hash of all transactions in the block
func getTransactionsHash(transactions []*Transaction) []byte {
	var transactionsData []byte
	for _, transaction := range transactions {
		transactionsData = append(transactionsData, transaction.Data...)
	}
	hash := sha256.New()
	hash.Write(transactionsData)
	return hash.Sum(nil)
}

// AddBlock adds a new block with the given transactions to the blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(prevBlock.Hash, transactions)
	bc.blocks = append(bc.blocks, newBlock)
}

// PrintBlockchain prints the information of all blocks in the blockchain
func (bc *Blockchain) PrintBlockchain() {
	for _, block := range bc.blocks {
		block.PrintBlock()
	}
}

// PrintBlock prints the information of a block
func (b *Block) PrintBlock() {
	fmt.Printf("Timestamp: %d\n", b.Timestamp)
	fmt.Printf("PrevBlockHash: %x\n", b.PrevBlockHash)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("Transactions: %v\n", b.Transactions)
	fmt.Println("--------------------")
}

// HandleP2PMessage handles incoming P2P messages
func HandleP2PMessage() {
	for {
		select {
		case newTransactions := <-p2p.Blocks:
			validateAndAddTransactions(newTransactions)
		}
	}
}

// validateAndAddTransactions validates and adds new transactions to the blockchain
func validateAndAddTransactions(transactions []*Transaction) {
	// Perform validation logic here if needed

	// Add valid transactions to the blockchain
	blockchain.AddBlock(transactions)
}