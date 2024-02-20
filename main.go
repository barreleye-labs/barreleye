package main

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/config"
	"github.com/barreleye-labs/barreleye/core/types"
	"log"
	"time"

	"github.com/barreleye-labs/barreleye/node"
)

func main() {
	nodeName := common.GetFlag()

	/* create hex key
	nodePrivateKey := types.GeneratePrivateKey()
	hexPrivateKey := hex.EncodeToString(crypto.FromECDSA(nodePrivateKey.Key))
	privateKey1, _ := crypto.HexToECDSA(hexPrivateKey)
	nodePrivateKey.Key = privateKey1
	fmt.Println("pk: ", hexPrivateKey)
	*/

	conf := config.GetConfig(nodeName)
	privateKey, err := types.CreatePrivateKey(conf.PrivateKey)
	if err != nil {
		panic("failed to create private key")
	}

	if nodeName == "genesis-node" {
		node1 := createNode("GENESIS-NODE", privateKey, ":3000", []string{":4000"}, ":9000")
		node1.Start()
	} else if nodeName == "nayoung" {
		node2 := createNode("NAYOUNG", privateKey, ":4000", []string{":3000"}, ":9001")
		node2.Start()
	} else if nodeName == "youngmin" {
		node3 := createNode("YOUNGMIN", privateKey, ":5000", []string{":4000"}, ":9002")
		node3.Start()
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

func createNode(id string, pk *types.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *node.Node {
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
