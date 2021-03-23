package main

type Node struct {
	identifier       string
	peerLocations    []string
	currentTerm      int
	msSinceHeartbeat int
	voteCount        int
	isLeader         bool
}

func (*Node) SendState(nothing string, reply *string) error {
	*reply = "YOLO"
	return nil
}
