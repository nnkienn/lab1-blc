package blockchain

import (
	"crypto/sha256"
	"fmt"
	"sync"
	"time"
)

type Transaction struct {
	Data []byte
}

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
}

type Blockchain struct {
	Blocks []*Block
	mutex  sync.Mutex
}

var blockchainInstance = &Blockchain{
	Blocks: []*Block{},
}

func GetBlockchainInstance() *Blockchain {
	return blockchainInstance
}

func (b *Block) SetHash() {
	headers := append(b.PrevBlockHash, getTransactionsHash(b.Transactions)...)
	hash := sha256.New()
	hash.Write(headers)
	b.Hash = hash.Sum(nil)
}

func getTransactionsHash(transactions []*Transaction) []byte {
	var transactionsData []byte
	for _, transaction := range transactions {
		transactionsData = append(transactionsData, transaction.Data...)
	}
	hash := sha256.New()
	hash.Write(transactionsData)
	return hash.Sum(nil)
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	var prevBlockHash []byte
	if len(bc.Blocks) > 0 {
		prevBlockHash = bc.Blocks[len(bc.Blocks)-1].Hash
	}

	newBlock := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
	}
	newBlock.SetHash()

	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) GetTransactions() []*Transaction {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if len(bc.Blocks) > 0 {
		return bc.Blocks[len(bc.Blocks)-1].Transactions
	}
	return nil
}

func (bc *Blockchain) PrintBlockchain() {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	for _, block := range bc.Blocks {
		block.PrintBlock()
	}
}

func (b *Block) PrintBlock() {
	fmt.Printf("Timestamp: %d\n", b.Timestamp)
	fmt.Printf("PrevBlockHash: %x\n", b.PrevBlockHash)
	fmt.Printf("Hash: %x\n", b.Hash)
	fmt.Printf("Transactions: %v\n", b.Transactions)
	fmt.Println("--------------------")
}
