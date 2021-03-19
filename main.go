package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"time"
)

func (node *Node) startServer() {
	rpc.Register(node)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf(":%v", node.portNumber))
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
	portNumber := os.Args[1]
	peerLocations := os.Args[2:]
	node := Node{portNumber, peerLocations, 0, 0, 0, false}

	go node.startServer()
	go node.runTerms()

	for {
		time.Sleep(time.Second * 100)
	}
}
