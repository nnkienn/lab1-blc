package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "github.com/gorilla/mux"
    "github.com/nnkienn/lab1-blc/merkletree"
    "github.com/nnkienn/lab1-blc/blockchain"
    "encoding/json"
)

var blockchain = blockchain.NewBlockchain()

// Hàm main thiết lập các đường dẫn và khởi động máy chủ
func main() {
    r := mux.NewRouter()

    // Các địa chỉ liên quan đến khối
    r.HandleFunc("/blocks", GetBlocksEndpoint).Methods("GET")
    r.HandleFunc("/mineBlock", MineBlockEndpoint).Methods("POST")

    // Địa chỉ WebSocket cho P2P đã được loại bỏ

    // Periodically broadcast the MerkleTree
    go func() {
        for {
            blockchain.BroadcastMerkleTree()
            time.Sleep(10 * time.Second)
        }
    }()

    port := ":3001"
    fmt.Println("Listening on port", port)
    log.Fatal(http.ListenAndServe(port, r))
}

// GetBlocksEndpoint log một thông điệp và in thông tin blockchain
func GetBlocksEndpoint(w http.ResponseWriter, r *http.Request) {
    fmt.Println("In thông tin blockchain:")
    blockchain.PrintBlockchain()
}

// MineBlockEndpoint đào một khối mới với các giao dịch cho trước
func MineBlockEndpoint(w http.ResponseWriter, r *http.Request) {
    var data map[string]string

    err := json.NewDecoder(r.Body).Decode(&data)
    if err != nil {
        http.Error(w, "Yêu cầu không hợp lệ", http.StatusBadRequest)
        return
    }

    blockchain.AddBlock([]*blockchain.Transaction{&blockchain.Transaction{Data: []byte(data["data"])}})
    blockchain.BroadcastMerkleTree()

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Khối được đào và được phát sóng"))
}
