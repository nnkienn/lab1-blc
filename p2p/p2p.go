// p2p.go

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
	Blocks chan []*blockchain.Block
}

var p2pInstance = &P2P{
	Nodes:  []*websocket.Conn{},
	Blocks: make(chan []*blockchain.Block),
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
	latestBlocks := blockchain.NewBlockchain().GetBlocks()
	p.Blocks <- latestBlocks
}

// ... (other methods)


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
	p.BroadcastBlockchain()

	go p.HandleP2PMessage(conn)
}

func (p *P2P) HandleP2PMessage(conn *websocket.Conn) {
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

func (p *P2P) BroadcastMessage(msg map[string]interface{}) {
	for _, node := range p.Nodes {
		if err := node.WriteJSON(msg); err != nil {
			log.Println(err)
		}
	}
}
