package main

type Node struct {
	portNumber       string
	peerLocations    []string
	currentTerm      int
	msSinceHeartbeat int
	voteCount        int
	isLeader         bool
}
