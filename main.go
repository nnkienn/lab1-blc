// main.go

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

	// Check if there are command line arguments
	if len(os.Args) > 1 {
		ports = os.Args[1:]
	} else {
		// If no arguments, use the default port 3001
		ports = []string{"3001"}
	}

	// Create a list of routers and worker goroutines
	routers := make([]*mux.Router, len(ports))
	for i, port := range ports {
		// Create a new router for each port
		routers[i] = mux.NewRouter()

		// Register router endpoints
		routers[i].HandleFunc("/blocks", GetBlocksEndpoint).Methods("GET")
		routers[i].HandleFunc("/mineBlock", MineBlockEndpoint).Methods("POST")
		routers[i].HandleFunc("/ws", p2pHandler.HandleP2PConnection)

		// Start a worker goroutine for each port
		go func(p string, r *mux.Router) {
			port := ":" + p
			fmt.Println("Listening on port", port)
			log.Fatal(http.ListenAndServe(port, r))
		}(port, routers[i])
	}

	go PeriodicallyBroadcastMerkleTree()

	// Wait for all worker goroutines to complete
	select {}
}

// GetBlocksEndpoint returns the blockchain as JSON
func GetBlocksEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	blocks := bc.GetBlocks()

	for _, block := range blocks {
		for _, transaction := range block.Transactions {
			fmt.Printf("Transaction Data: %s\n", string(transaction.Data))
		}
	}

	json.NewEncoder(w).Encode(blocks)
}

// MineBlockEndpoint mines a new block and broadcasts it
func MineBlockEndpoint(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mine a new block with the provided data
	p2pHandler.MineAndBroadcastBlock(data["data"])

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Block mined and broadcasted"))
}

// PeriodicallyBroadcastMerkleTree broadcasts the Merkle tree periodically
func PeriodicallyBroadcastMerkleTree() {
	for {
		p2pHandler.BroadcastMerkleTree()
		time.Sleep(10 * time.Second)
	}
}
