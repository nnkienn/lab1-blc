package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/nnkienn/lab1-blc/blockchain"
	"github.com/nnkienn/lab1-blc/p2p"
)

var bc = blockchain.NewBlockchain()
var p2pHandler = p2p.GetP2PInstance()

func main() {
	var ports []string

	// Kiểm tra xem có đối số dòng lệnh không
	if len(os.Args) > 1 {
		ports = os.Args[1:]
	} else {
		// Nếu không có đối số, sử dụng cổng mặc định 3001
		ports = []string{"3001"}
	}

	// Tạo một danh sách các router và các worker goroutine
	routers := make([]*mux.Router, len(ports))
	for i, port := range ports {
		// Tạo một router mới cho mỗi cổng
		routers[i] = mux.NewRouter()

		// Đăng ký các endpoint của router
		routers[i].HandleFunc("/blocks", GetBlocksEndpoint).Methods("GET")
		routers[i].HandleFunc("/mineBlock", MineBlockEndpoint).Methods("POST")
		routers[i].HandleFunc("/ws", p2pHandler.HandleP2PConnection)

		// Bắt đầu một worker goroutine cho mỗi cổng
		go func(p string, r *mux.Router) {
			port := ":" + p
			fmt.Println("Listening on port", port)
			log.Fatal(http.ListenAndServe(port, r))
		}(port, routers[i])
	}

	go PeriodicallyBroadcastMerkleTree()

	// Đợi tất cả các worker goroutine hoàn thành
	select {}
}

// ... (rest of your code)


func GetBlocksEndpoint(w http.ResponseWriter, r *http.Request) {
	// Return the blockchain as JSON
	w.Header().Set("Content-Type", "application/json")

	blocks := bc.GetBlocks()

	// Loop through transactions and print their data
	for _, block := range blocks {
		for _, transaction := range block.Transactions {
			fmt.Printf("Transaction Data: %s\n", string(transaction.Data))
		}
	}

	json.NewEncoder(w).Encode(blocks)
}

func MineBlockEndpoint(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction := &blockchain.Transaction{Data: []byte(data["data"])}
	bc.AddBlock([]*blockchain.Transaction{transaction})
	p2pHandler.BroadcastBlockchain()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Block mined and broadcasted"))
}

func PeriodicallyBroadcastMerkleTree() {
	for {
		p2pHandler.BroadcastBlockchain()
		time.Sleep(10 * time.Second)
	}
}
