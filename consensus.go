package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type ConsensusModule struct {
	nextIndex   map[string]int
	matchIndex  map[string]int
	commitIndex int
}

type LogEntry struct {
	command    map[string]string
	termNumber int
}

type Log []LogEntry

var logMutex sync.Mutex

func (cm *ConsensusModule) startServer(nodeState *Node) {

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		var key string
		var value string
		key = strings.Split(r.URL.RawQuery, "=")[0]
		value = strings.Split(r.URL.RawQuery, "=")[1]
		fmt.Println("Received the set command: ", key, value)

		nodeState.log = append(nodeState.log, LogEntry{map[string]string{key: value}, nodeState.currentTerm})
	})

	http.ListenAndServe(":8080", nil)
}
