package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
)

type Node struct {
	id         int
	left_chan  <-chan Token
	right_chan chan Token
}

type Token struct {
	Data string `json:"data"`
	Recv int    `json:"recv"`
	Ttl  int    `json:"ttl"`
}

func (node *Node) Run() {
	msg := <-node.left_chan
	switch {
	case msg.Recv == node.id:
		fmt.Println("NodeID", node.id, "Recieve ( Recv:", msg.Recv, "Data:", msg.Data, " Ttl:", msg.Ttl, ")")
	case msg.Ttl > 0:
		msg.Ttl -= 1
		fmt.Println("NodeID", node.id, "Pass ( Recv:", msg.Recv, " Data:", msg.Data, " Ttl:", msg.Ttl, ")")
		node.right_chan <- msg
	default:
		fmt.Println("NodeID", node.id, "Expired ( Recv:", msg.Recv, "Data:", msg.Data, " Ttl:", msg.Ttl, ")")
	}
}

func initialize(n int) []*Node {
	ring := make([]*Node, 0, n)

	ring = append(ring, &Node{id: 0, right_chan: make(chan Token)})

	for i := 1; i < n; i++ {
		ring = append(ring, &Node{id: i, left_chan: ring[i-1].right_chan, right_chan: make(chan Token)})
	}

	ring[0].left_chan = ring[n-1].right_chan

	return ring
}

func sendMsg(w http.ResponseWriter, r *http.Request) {
	token := Token{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&token)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Println("Origin data: ", token.Data, " ", token.Recv, "", token.Ttl)

	for _, node := range ring {
		go node.Run()
	}

	ring[len(ring)-1].right_chan <- token
}

var ring []*Node

func main() {

	var nodeCnt = flag.Int("numOfNodes", 3, "Number of nodes")
	flag.Parse()

	fmt.Println("Server starts with number of nodes = ", *nodeCnt)
	ring = initialize(*nodeCnt)

	http.HandleFunc("/", sendMsg)
	http.ListenAndServe(":3000", nil)
}
