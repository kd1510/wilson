package main

import (
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

type Vote struct {
	GivenVote bool
}

//We need this because a heartbeat could be received while reading the timer
var heartBeatMutex sync.Mutex
var voteMutex sync.Mutex

//RPC Methods
func (node *Node) ReceiveHeartbeat(master string, reply *string) error {
	fmt.Printf("Node %v received a heartbeat from master %v\n", node.portNumber, master)
	heartBeatMutex.Lock()
	node.msSinceHeartbeat = 0
	heartBeatMutex.Unlock()
	*reply = "thanks mate"
	return nil
}

func (node *Node) SendVote(candidate string, reply *Vote) error {
	if node.isLeader != true {
		*reply = Vote{true}
	} else {
		*reply = Vote{false}
	}
	return nil
}

func (node *Node) sendHeartBeat() {
	for _, peerLoc := range node.peerLocations {
		go func(peerLoc string) {
			fmt.Printf("Node %v sending heartbeat to %v\n", node.portNumber, peerLoc)
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%v", peerLoc))

			if err != nil {
				fmt.Printf("Follower %v must be down, couldn't connect\n", peerLoc)
				return
			}
			defer client.Close()

			var response string
			err = client.Call("Node.ReceiveHeartbeat", node.portNumber, &response)
			if err != nil {
				fmt.Printf("Error sending heartbeat to follower %v\n", peerLoc)
				return
			}
		}(peerLoc)
	}
}

func (node *Node) heartbeatTimeout(cutoff int) {
	fmt.Printf("Timeout %v\n", cutoff)
	for {
		heartBeatMutex.Lock()
		if node.msSinceHeartbeat > cutoff {
			heartBeatMutex.Unlock()
			node.initiateElection()
			return
		} else {
			time.Sleep(1 * time.Millisecond)
			node.msSinceHeartbeat++
			heartBeatMutex.Unlock()
		}
	}
}

func (node *Node) requestVotes() {
	var wg sync.WaitGroup
	for _, peerLoc := range node.peerLocations {
		wg.Add(1)
		go func(peerLoc string) {
			fmt.Printf("Node %v requesting vote from %v\n", node.portNumber, peerLoc)
			client, err := rpc.DialHTTP("tcp", fmt.Sprintf("localhost:%v", peerLoc))
			if err != nil {
				fmt.Printf("Couldn't connect to peer %v\n", peerLoc)
				wg.Done()
				return
			}

			defer client.Close()
			var vote Vote
			_ = client.Call("Node.SendVote", node.portNumber, &vote)
			if vote.GivenVote == true {
				voteMutex.Lock()
				node.voteCount++
				voteMutex.Unlock()
			}
			wg.Done()
		}(peerLoc)
	}
	wg.Wait()
	fmt.Printf("Finished voting round, recieved %v votes\n", node.voteCount)
	return
}

func (node *Node) initiateElection() {
	fmt.Printf("%v Initiating election as reached timeout cutoff\n", node.portNumber)
	//If a node initiates, votes for itself
	voteMutex.Lock()
	node.voteCount++
	voteMutex.Unlock()

	//Request votes from peers
	node.requestVotes()

	//If receives majority votes then becomes leader
	voteMutex.Lock()
	if node.voteCount >= (len(node.peerLocations)+1)-(len(node.peerLocations)/2) {
		node.isLeader = true
		node.voteCount = 0
	} else {
		node.voteCount = 0
	}
	voteMutex.Unlock()
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
