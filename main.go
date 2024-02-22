package main

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"log"
	"time"

	"github.com/barreleye-labs/barreleye/node"
)

func main() {
	common.ParseFlag()
	nodeName := common.GetFlag("name")
	port := common.GetFlag("port")
	peer := common.GetFlag("peer")
	httpPort := common.GetFlag("http.port")
	key := common.GetFlag("key")

	/* create hex key
	nodePrivateKey := types.GeneratePrivateKey()
	hexPrivateKey := hex.EncodeToString(crypto.FromECDSA(nodePrivateKey.Key))
	privateKey1, _ := crypto.HexToECDSA(hexPrivateKey)
	nodePrivateKey.Key = privateKey1
	fmt.Println("pk: ", hexPrivateKey)
	*/

	privateKey, err := types.CreatePrivateKey(key)
	if err != nil {
		panic("failed to create private key")
	}

	node1 := createNode(nodeName, privateKey, ":"+port, []string{peer}, ":"+httpPort)
	node1.Start()

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
