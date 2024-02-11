package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/barreleye-labs/barreleye/node"
)

func main() {
	fmt.Println("start")

	var nodeName string = ""
	flag.StringVar(&nodeName, "nodeName", "", "Node name")
	flag.Parse()

	// file, _ := os.Open("config/config.json")
	// defer file.Close()
	// decoder := json.NewDecoder(file)
	// nodeInfo := config.NodeInfo{}
	// err := decoder.Decode(&nodeInfo)
	// if err != nil {
	// fmt.Println("error:", err)
	// }

	if nodeName == "node1" {
		validatorPrivKey := crypto.GeneratePrivateKey()
		node := createNode("NODE1", &validatorPrivKey, ":3000", []string{":4000"}, ":9000")
		node.Start()
	} else if nodeName == "node2" {
		node := createNode("NODE2", nil, ":4000", []string{":3000"}, ":9001")
		node.Start()
	} else if nodeName == "node3" {
		node := createNode("NODE3", nil, ":5000", []string{":3000", ":4000"}, "")
		node.Start()
	}

	// fmt.Println("kyma:", nodeInfo.Node1.Endpoint)

	// validatorPrivKey := crypto.GeneratePrivateKey()
	// localNode := makeServer("LOCAL_NODE", &validatorPrivKey, "localhost:3000", []string{"localhost:4000"}, ":9000")
	// go localNode.Start()

	// remoteNode := makeServer("REMOTE_NODE", nil, "localhost:4000", []string{"localhost:5000"}, "")
	// go remoteNode.Start()

	// remoteNodeB := makeServer("REMOTE_NODE_B", nil, "localhost:5000", nil, "")
	// go remoteNodeB.Start()

	// go func() {
	// 	time.Sleep(11 * time.Second)

	// 	lateNode := makeServer("LATE_NODE", nil, ":6000", []string{"localhost:4000"}, "")
	// 	go lateNode.Start()
	// }()

	time.Sleep(1 * time.Second)

	select {}
}

func createNode(id string, pk *crypto.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *node.Node {
	opts := node.NodeOpts{
		APIListenAddr: apiListenAddr,
		SeedNodes:     seedNodes,
		ListenAddr:    addr,
		PrivateKey:    pk,
		ID:            id,
	}

	s, err := node.NewNode(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
