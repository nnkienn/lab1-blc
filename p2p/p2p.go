package p2p

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nnkienn/lab1-blc/blockchain"

)

// P2P represents the Peer-to-Peer communication
type P2P struct {
	Nodes  []*websocket.Conn
	Mutex  sync.Mutex
	Blocks chan []*Transaction
}

var p2p = &P2P{
	Nodes:  []*websocket.Conn{},
	Blocks: make(chan []*Transaction),
}

// RegisterNode registers a new node
func (p *P2P) RegisterNode(conn *websocket.Conn) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.Nodes = append(p.Nodes, conn)
}

// BroadcastMerkleTree broadcasts the MerkleTree data to all nodes
func (p *P2P) BroadcastMerkleTree() {
	p.Blocks <- merkletree.GetMerkleTreeData()
}

// HandleP2PConnection handles a new P2P connection
func (p *P2P) HandleP2PConnection(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	p.RegisterNode(conn)

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}

		if msg["type"] == "merkle" {
			p.BroadcastMerkleTree()
		}
	}
}
