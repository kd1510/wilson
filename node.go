package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	identifier       string
	peerLocations    []string
	log              Log
	currentTerm      int
	msSinceHeartbeat int
	voteCount        int
	isLeader         bool
}

func (n *Node) appendEntry(entry LogEntry, currentLogIndex *int) error {
	n.log = append(n.log, entry)
	*currentLogIndex = len(n.log)
	return nil
}

func (node *Node) runTerms() {

	go func() {
		for {
			fmt.Println("Current log: ", node.log)
			time.Sleep(3 * time.Second)
		}
	}()

	for {

		if node.isLeader == true {

			fmt.Println("Elected Leader")

			//Start sending out heartbeat to all followers
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				for {
					node.sendHeartBeat()
					time.Sleep(time.Millisecond * 500)
				}
			}()

			//Listen to commands from clients
			//replicate to followers
			go func() {
				fmt.Println("Opening HTTP Server!")
				cm := new(ConsensusModule)
				cm.startServer(node)
			}()

			//Periodic log state
			wg.Wait()

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
