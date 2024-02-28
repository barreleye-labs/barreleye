package main

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"log"
	"strings"
	"time"

	"github.com/barreleye-labs/barreleye/node"
)

func main() {
	common.ParseFlag()
	nodeName := common.GetFlag("name")
	port := common.GetFlag("port")
	peers := common.GetFlag("peers")
	httpPort := common.GetFlag("http.port")
	key := common.GetFlag("key")

	peerArr := []string{}
	if peers != "none" {
		peerArr = strings.Split(peers, ",")
	}

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

	n := createNode(nodeName, privateKey, ":"+port, peerArr, ":"+httpPort)
	n.Start()

	time.Sleep(1 * time.Second)

	select {}
}

func createNode(id string, pk *types.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *node.Node {
	opts := node.NodeOpts{
		APIListenAddr: apiListenAddr,
		SeedNodes:     seedNodes,
		ListenAddr:    addr,
		PrivateKey:    pk,
		Name:          id,
	}

	s, err := node.NewNode(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}
