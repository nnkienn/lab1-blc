package p2p

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/nnkienn/lab1-blc/blockchain"
)

type P2P struct {
	Nodes  []*websocket.Conn
	Mutex  sync.Mutex
	Blocks chan []*blockchain.Transaction
}

var p2pInstance = &P2P{
	Nodes:  []*websocket.Conn{},
	Blocks: make(chan []*blockchain.Transaction),
}

func GetP2PInstance() *P2P {
	return p2pInstance
}

func (p *P2P) RegisterNode(conn *websocket.Conn) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	p.Nodes = append(p.Nodes, conn)
}

func (p *P2P) BroadcastBlockchain() {
	latestTransactions := blockchain.GetBlockchainInstance().GetTransactions()
	p.Blocks <- latestTransactions
}

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

		if msg["type"] == "blocks" {
			p.BroadcastBlockchain()
		}
	}
}
