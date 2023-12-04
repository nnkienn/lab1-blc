package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/nnkienn/lab1-blc/blockchain"
)

// ... (rest of the code)

var (
	mu           sync.Mutex
	blockchain   []*blockchain.Block
	transactions []*blockchain.Transaction
	peers        []string
)

func main() {
	// Run multiple instances on different ports
	go runBlockchain(":8080")
	go runBlockchain(":8081")

	// Allow time for instances to start
	time.Sleep(2 * time.Second)

	// Connect peers
	connectPeers()

	select {}
}

func runBlockchain(port string) {
	r := mux.NewRouter()

	r.HandleFunc("/blocks", GetBlocks).Methods("GET")
	r.HandleFunc("/addblock", AddBlock).Methods("POST")
	r.HandleFunc("/addtransaction", AddTransaction).Methods("POST")

	http.Handle("/", r)

	fmt.Printf("Server listening on %s\n", port)
	http.ListenAndServe(port, nil)
}

func connectPeers() {
	mu.Lock()
	defer mu.Unlock()

	for _, port := range []string{":8080", ":8081"} {
		if len(peers) == 0 || port != peers[0] {
			peers = append(peers, fmt.Sprintf("http://localhost%s", port))
		}
	}
}

func GetBlocks(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	json.NewEncoder(w).Encode(blockchain)
}

func AddBlock(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	prevBlock := blockchain[len(blockchain)-1]
	newBlock := blockchain.NewBlock(transactions, prevBlock.Hash)
	transactions = nil
	blockchain = append(blockchain, newBlock)

	// Broadcast the new block to peers
	go broadcastBlock(newBlock)

	json.NewEncoder(w).Encode(newBlock)
}

func broadcastBlock(newBlock *blockchain.Block) {
	mu.Lock()
	defer mu.Unlock()

	for _, peer := range peers {
		go func(peerAddr string) {
			resp, err := http.Post(peerAddr+"/addblock", "application/json", nil)
			if err != nil {
				fmt.Println("Error broadcasting block to", peerAddr, ":", err)
				return
			}
			defer resp.Body.Close()
		}(peer)
	}
}

func AddTransaction(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	newTransaction := blockchain.NewTransaction([]byte(data.Data))
	mu.Lock()
	defer mu.Unlock()
	transactions = append(transactions, newTransaction)

	w.WriteHeader(http.StatusCreated)
}
