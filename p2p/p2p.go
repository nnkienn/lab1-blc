package p2p

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"github.com/nnkienn/lab1-blc/blockchain"
)

// P2PNode represents a P2P node
type P2PNode struct {
	Port       string
	Blockchain *blockchain.Blockchain
}

// StartP2PServer starts the P2P server
func (node *P2PNode) StartP2PServer() {
	ln, err := net.Listen("tcp", ":"+node.Port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	fmt.Println("P2P Server is listening on port", node.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go node.handleConnection(conn)
	}
}

// handleConnection handles incoming connections
func (node *P2PNode) handleConnection(conn net.Conn) {
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Text()
		fmt.Println("Received:", data)

		switch {
		case strings.HasPrefix(data, "QUERY"):
			queryParts := strings.Split(data, ":")
			if len(queryParts) == 2 {
				blockHash := queryParts[1]
				node.queryBlock(blockHash, conn)
			} else {
				node.sendResponse(conn, "Invalid QUERY format")
			}
		case strings.HasPrefix(data, "ADD"):
			dataParts := strings.Split(data, ":")
			if len(dataParts) == 2 {
				transactionData := dataParts[1]
				node.addTransaction([]byte(transactionData))
				node.sendResponse(conn, "Transaction added to the blockchain")
			} else {
				node.sendResponse(conn, "Invalid ADD format")
			}
		default:
			node.sendResponse(conn, "Unknown command")
		}
	}
}

// queryBlock handles the QUERY command
func (node *P2PNode) queryBlock(blockHash string, conn net.Conn) {
	node.Blockchain.Mu.Lock()
	defer node.Blockchain.Mu.Unlock()

	for _, block := range node.Blockchain.Blocks {
		if block.HexHash() == blockHash {
			response, err := json.Marshal(block)
			if err != nil {
				node.sendResponse(conn, "Error marshalling block")
				return
			}
			node.sendResponse(conn, string(response))
			return
		}
	}

	node.sendResponse(conn, "Block not found")
}

// addTransaction adds a new transaction to the blockchain
func (node *P2PNode) addTransaction(data []byte) {
	transaction := blockchain.NewTransaction(data)
	node.Blockchain.Mu.Lock()
	defer node.Blockchain.Mu.Unlock()

	// Placeholder: In a real P2P network, transactions would be broadcasted to other nodes
	// Here, we simply add the transaction to the local blockchain
	block := blockchain.NewBlock([]*blockchain.Transaction{transaction}, node.Blockchain.Blocks[len(node.Blockchain.Blocks)-1].Hash)
	node.Blockchain.Blocks = append(node.Blockchain.Blocks, block)
}

// sendResponse sends a response to the client
func (node *P2PNode) sendResponse(conn net.Conn, response string) {
	conn.Write([]byte(response + "\n"))
}
