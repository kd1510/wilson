package main

import (
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"time"
)

func (node *Node) startServer() {
	rpc.Register(node)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "0.0.0.0:12345")
	if e != nil {
		panic("listen error:")
	}
	go http.Serve(l, nil)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	identifier := os.Getenv("IDENTIFIER")
	peerLocations := strings.Split(os.Getenv("PEERS"), ",")
	node := Node{identifier, peerLocations, make(Log, 0), 0, 0, 0, false}

	go node.startServer()

	//Can heartbeat/elections run in different thread to the log replication?
	go node.runTerms()

	for {
		time.Sleep(time.Second * 100)
	}
}
