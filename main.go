package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/common/util"
	"github.com/barreleye-labs/barreleye/core/types"
	"log"
	"net/http"
	"time"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/barreleye-labs/barreleye/network"
)

func main() {
	fmt.Println("start")

	var nodeName string = ""
	flag.StringVar(&nodeName, "nodeName", "", "Node name")
	flag.Parse()
	fmt.Println("nono: ", nodeName)

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
		node := makeServer("NODE1", &validatorPrivKey, ":3000", []string{":4000"}, ":9000")
		node.Start()
	} else if nodeName == "node2" {
		node := makeServer("NODE2", nil, ":4000", []string{":3000"}, "")
		node.Start()
	} else if nodeName == "node3" {
		node := makeServer("NODE3", nil, ":5000", []string{":4000"}, "")
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

	// if err := sendTransaction(validatorPrivKey); err != nil {
	// 	panic(err)
	// }

	// collectionOwnerPrivKey := crypto.GeneratePrivateKey()
	// collectionHash := createCollectionTx(collectionOwnerPrivKey)

	// txSendTicker := time.NewTicker(1 * time.Second)
	// go func() {
	// 	for i := 0; i < 20; i++ {
	// 		nftMinter(collectionOwnerPrivKey, collectionHash)

	// 		<-txSendTicker.C
	// 	}
	// }()

	select {}
}

func sendTransaction(privKey crypto.PrivateKey) error {
	toPrivKey := crypto.GeneratePrivateKey()

	tx := types.NewTransaction(nil)
	tx.To = toPrivKey.PublicKey()
	tx.Value = 666

	if err := tx.Sign(privKey); err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)

	return err
}

func makeServer(id string, pk *crypto.PrivateKey, addr string, seedNodes []string, apiListenAddr string) *network.Server {
	opts := network.ServerOpts{
		APIListenAddr: apiListenAddr,
		SeedNodes:     seedNodes,
		ListenAddr:    addr,
		PrivateKey:    pk,
		ID:            id,
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

func createCollectionTx(privKey crypto.PrivateKey) common.Hash {
	tx := types.NewTransaction(nil)
	tx.TxInner = types.CollectionTx{
		Fee:      200,
		MetaData: []byte("chicken and egg collection!"),
	}
	tx.Sign(privKey)

	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	return tx.GetHash(types.TxHasher{})
}

func nftMinter(privKey crypto.PrivateKey, collection common.Hash) {
	metaData := map[string]any{
		"power":  8,
		"health": 100,
		"color":  "green",
		"rare":   "yes",
	}

	metaBuf := new(bytes.Buffer)
	if err := json.NewEncoder(metaBuf).Encode(metaData); err != nil {
		panic(err)
	}

	tx := types.NewTransaction(nil)
	tx.TxInner = types.MintTx{
		Fee:             200,
		NFT:             util.RandomHash(),
		MetaData:        metaBuf.Bytes(),
		Collection:      collection,
		CollectionOwner: privKey.PublicKey(),
	}
	tx.Sign(privKey)

	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", "http://localhost:9000/tx", buf)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
