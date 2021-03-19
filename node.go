package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Node struct {
	portNumber       string
	peerLocations    []string
	currentTerm      int
	msSinceHeartbeat int
	voteCount        int
	isLeader         bool
}

func (node *Node) runTerms() {
	for {
		fmt.Printf("Term: %v |||  Is Leader: %v ||| \n", node.currentTerm, node.isLeader)

		if node.isLeader == true {

			//Start sending out heartbeat to all followers
			node.sendHeartBeat()
			time.Sleep(time.Millisecond * 500)

		} else {

			//Check if we need to initiate a leader election
			rand.Seed(time.Now().UTC().UnixNano())
			cutoff := rand.Intn(1000) + 2000 //from 150 to 300ms
			node.heartbeatTimeout(cutoff)
			node.currentTerm++

			heartBeatMutex.Lock()
			node.msSinceHeartbeat = 0
			heartBeatMutex.Unlock()
		}
	}
}
