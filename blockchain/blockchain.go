package blockchain

import (
    "crypto/sha256"
    "fmt"
    "time"
)

// Transaction là một giao dịch cơ bản trong blockchain
type Transaction struct {
    Data []byte
}

// Block đại diện cho một khối trong blockchain
type Block struct {
    Timestamp     int64
    Transactions  []*Transaction
    PrevBlockHash []byte
    Hash          []byte
}

// Blockchain đại diện cho toàn bộ blockchain
type Blockchain struct {
    Blocks []*Block
}

// NewBlockchain tạo một blockchain mới với khối khởi tạo
func NewBlockchain() *Blockchain {
    genesisTransaction := &Transaction{Data: []byte("Giao dịch khởi tạo")}
    genesisBlock := NewBlock([]byte{}, []*Transaction{genesisTransaction})
    return &Blockchain{Blocks: []*Block{genesisBlock}}
}

// NewBlock tạo một khối mới với các giao dịch và hash khối trước đó
func NewBlock(prevBlockHash []byte, transactions []*Transaction) *Block {
    block := &Block{
        Timestamp:     getCurrentTimestamp(),
        Transactions:  transactions,
        PrevBlockHash: prevBlockHash,
    }
    block.SetHash()
    return block
}

// getCurrentTimestamp trả về timestamp hiện tại
func getCurrentTimestamp() int64 {
    return time.Now().Unix()
}

// SetHash tính toán và đặt hash cho khối dựa trên các trường của nó
func (b *Block) SetHash() {
    headers := append(b.PrevBlockHash, getTransactionsHash(b.Transactions)...)
    hash := sha256.New()
    hash.Write(headers)
    b.Hash = hash.Sum(nil)
}

// getTransactionsHash tính toán hash của tất cả các giao dịch trong khối
func getTransactionsHash(transactions []*Transaction) []byte {
    var transactionsData []byte
    for _, transaction := range transactions {
        transactionsData = append(transactionsData, transaction.Data...)
    }
    hash := sha256.New()
    hash.Write(transactionsData)
    return hash.Sum(nil)
}

// AddBlock thêm một khối mới với các giao dịch cho trước vào blockchain
func (bc *Blockchain) AddBlock(transactions []*Transaction) {
    prevBlock := bc.Blocks[len(bc.Blocks)-1]
    newBlock := NewBlock(prevBlock.Hash, transactions)
    bc.Blocks = append(bc.Blocks, newBlock)
}

// PrintBlockchain in thông tin của tất cả các khối trong blockchain
func (bc *Blockchain) PrintBlockchain() {
    for _, block := range bc.Blocks {
        block.PrintBlock()
    }
}

// PrintBlock in thông tin của một khối
func (b *Block) PrintBlock() {
    fmt.Printf("Timestamp: %d\n", b.Timestamp)
    fmt.Printf("PrevBlockHash: %x\n", b.PrevBlockHash)
    fmt.Printf("Hash: %x\n", b.Hash)
    fmt.Printf("Transactions: %v\n", b.Transactions)
    fmt.Println("--------------------")
}

// BroadcastMerkleTree log một thông điệp khi phát sóng dữ liệu của cây Merkle
func (bc *Blockchain) BroadcastMerkleTree() {
    fmt.Println("Phát sóng dữ liệu cây Merkle")
    // Gọi hàm GetMerkleTreeData() nếu nó tồn tại
}
