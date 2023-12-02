
// main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"merkletree/merkletree"
	"p2p/p2p"
	"blockchain/blockchain"
)

var blockchain = blockchain.NewBlockchain()

func main() {
	r := mux.NewRouter()

	// Block-related endpoints
	r.HandleFunc("/blocks", GetBlocksEndpoint).Methods("GET")
	r.HandleFunc("/mineBlock", MineBlockEndpoint).Methods("POST")

	// WebSocket endpoint for P2P
	r.HandleFunc("/ws", p2p.WebSocketHandler)

	// Start P2P message handler
	go blockchain.HandleP2PMessage()

	// Periodically broadcast the MerkleTree
	go func() {
		for {
			p2p.BroadcastMerkleTree()
			time.Sleep(10 * time.Second)
		}
	}()

	port := ":3001"
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}

// GetBlocksEndpoint returns the information of all blocks in the blockchain
func GetBlocksEndpoint(w http.ResponseWriter, r *http.Request) {
	blockchain.PrintBlockchain()
}

// MineBlockEndpoint mines a new block with the given transactions
func MineBlockEndpoint(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	blockchain.AddBlock([]*blockchain.Transaction{&blockchain.Transaction{Data: []byte(data["data"])}})
	p2p.BroadcastMerkleTree()

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Block mined and broadcasted"))
}