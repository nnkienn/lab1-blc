package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nnkienn/lab1-blc/blockchain"
	"github.com/nnkienn/lab1-blc/p2p"
)

var bc = blockchain.GetBlockchainInstance()
var p2pHandler = p2p.GetP2PInstance()

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/blocks", GetBlocksEndpoint).Methods("GET")
	r.HandleFunc("/mineBlock", MineBlockEndpoint).Methods("POST")

	r.HandleFunc("/ws", p2pHandler.HandleP2PConnection)

	go PeriodicallyBroadcastMerkleTree()

	port := ":3001"
	fmt.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(port, r))
}


func GetBlocksEndpoint(w http.ResponseWriter, r *http.Request) {
	bc.PrintBlockchain()
}

func MineBlockEndpoint(w http.ResponseWriter, r *http.Request) {
	var data map[string]string

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	bc.AddBlock([]*blockchain.Transaction{&blockchain.Transaction{Data: []byte(data["data"])}})
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
