package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"net"
	"os"
	"sync"
	"time"

	"github.com/barreleye-labs/barreleye/core"
	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/barreleye-labs/barreleye/restful"
	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type NodeOpts struct {
	APIListenAddr string
	SeedNodes     []string
	ListenAddr    string
	TCPTransport  *TCPTransport
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Node struct {
	TCPTransport *TCPTransport
	peerCh       chan *TCPPeer

	mu      sync.RWMutex
	peerMap map[net.Addr]*TCPPeer

	NodeOpts
	mempool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
	txChan      chan *types.Transaction
}

func NewNode(opts NodeOpts) (*Node, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "addr", opts.ID)
	}

	chain, err := core.NewBlockchain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	txChan := make(chan *types.Transaction)

	if len(opts.APIListenAddr) > 0 {
		apiNodeCfg := restful.ServerConfig{
			Logger:     opts.Logger,
			ListenAddr: opts.APIListenAddr,
		}
		apiNode := restful.NewServer(apiNodeCfg, chain, txChan)
		go apiNode.Start()

		opts.Logger.Log("msg", "JSON API Node running", "port", opts.APIListenAddr)
	}

	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddr, peerCh)

	s := &Node{
		TCPTransport: tr,
		peerCh:       peerCh,
		peerMap:      make(map[net.Addr]*TCPPeer),
		NodeOpts:     opts,
		chain:        chain,
		mempool:      NewTxPool(1000),
		isValidator:  opts.PrivateKey != nil,
		rpcCh:        make(chan RPC),
		quitCh:       make(chan struct{}, 1),
		txChan:       txChan,
	}

	s.TCPTransport.peerCh = peerCh

	// If we dont got any processor from the Node options, we going to use
	// the Node as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (n *Node) bootstrapNetwork() {
	for _, addr := range n.SeedNodes {
		fmt.Println("trying to connect to ", addr)

		go func(addr string) {
			conn, err := net.Dial("tcp", addr)
			if err != nil {
				fmt.Printf("could not connect to %+v\n", conn)
				return
			}

			n.peerCh <- &TCPPeer{
				conn: conn,
			}
		}(addr)
	}
}

func (n *Node) Start() {
	n.TCPTransport.Start()

	time.Sleep(time.Second * 1)

	n.bootstrapNetwork()

	n.Logger.Log("msg", "🤝 accepting TCP connection on", "addr", n.ListenAddr, "id", n.ID)

free:
	for {
		select {
		case peer := <-n.peerCh:
			n.peerMap[peer.conn.RemoteAddr()] = peer

			go peer.readLoop(n.rpcCh)

			if err := n.sendGetStatusMessage(peer); err != nil {
				n.Logger.Log("err", err)
				continue
			}

			n.Logger.Log("msg", "🙋 peer added to the Node", "outgoing", peer.Outgoing, "addr", peer.conn.RemoteAddr())

		case tx := <-n.txChan:
			if err := n.processTransaction(tx); err != nil {
				n.Logger.Log("process TX error", err)
			}

		case rpc := <-n.rpcCh:
			msg, err := n.RPCDecodeFunc(rpc)
			if err != nil {
				n.Logger.Log("RPC error", err)
				continue
			}

			if err := n.RPCProcessor.ProcessMessage(msg); err != nil {
				if err != core.ErrBlockKnown {
					n.Logger.Log("error", err)
				}
			}

		case <-n.quitCh:
			break free
		}
	}

	n.Logger.Log("msg", "Node is shutting down")
}

func (n *Node) validatorLoop() {
	ticker := time.NewTicker(n.BlockTime)

	n.Logger.Log("msg", "Starting validator loop", "blockTime", n.BlockTime)

	for {
		n.Logger.Log("msg", "🍀 creating new block")

		if err := n.createNewBlock(); err != nil {
			n.Logger.Log("create block error", err)
		}

		<-ticker.C
	}
}

func (n *Node) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *types.Transaction:
		return n.processTransaction(t)
	case *types.Block:
		return n.processBlock(t)
	case *GetStatusMessage:
		return n.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		return n.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		return n.processGetBlocksMessage(msg.From, t)
	case *BlocksMessage:
		return n.processBlocksMessage(msg.From, t)
	}

	return nil
}

func (n *Node) processGetBlocksMessage(from net.Addr, data *GetBlocksMessage) error {
	n.Logger.Log("msg", "📬 received getBlocks message", "from", from)

	var (
		blocks    = []*types.Block{}
		ourHeight = n.chain.Height()
	)

	if data.To == 0 {
		for i := int(data.From); i <= int(ourHeight); i++ {
			block, err := n.chain.GetBlock(uint32(i))
			if err != nil {
				return err
			}

			blocks = append(blocks, block)
		}
	}

	blocksMsg := &BlocksMessage{
		Blocks: blocks,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(blocksMsg); err != nil {
		return err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	msg := NewMessage(MessageTypeBlocks, buf.Bytes())
	peer, ok := n.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
	}

	return peer.Send(msg.Bytes())
}

func (n *Node) sendGetStatusMessage(peer *TCPPeer) error {
	var (
		getStatusMsg = new(GetStatusMessage)
		buf          = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
	return peer.Send(msg.Bytes())
}

func (n *Node) broadcast(payload []byte) error {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for netAddr, peer := range n.peerMap {
		if err := peer.Send(payload); err != nil {
			fmt.Printf("peer send error => addr %s [err: %s]\n", netAddr, err)
		}
	}
	return nil
}

func (n *Node) processBlocksMessage(from net.Addr, data *BlocksMessage) error {
	n.Logger.Log("msg", "📦 received BLOCKS", "from", from, "aff:", data.Blocks)

	for _, block := range data.Blocks {
		if err := n.chain.AddBlock(block); err != nil {
			n.Logger.Log("error", err.Error())
			return err
		}
	}

	return nil
}

func (n *Node) processStatusMessage(from net.Addr, data *StatusMessage) error {
	n.Logger.Log("msg", "📬 received STATUS message", "from", from)

	// 전달 받은 블록 높이보다 현재 나의 블록체인의 블록 높이가 같거나 클 경우.
	if data.CurrentHeight <= n.chain.Height() {
		n.Logger.Log("msg", "cannot sync blockHeight to low", "curHeight", n.chain.Height(), "theirHeight", data.CurrentHeight, "addr", from)
		return nil
	}

	go n.requestBlocksLoop(from)

	return nil
}

func (n *Node) processGetStatusMessage(from net.Addr, data *GetStatusMessage) error {
	n.Logger.Log("msg", "📬 received getStatus message", "from", from)

	StatusMessage := &StatusMessage{
		CurrentHeight: n.chain.Height(),
		ID:            n.ID,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(StatusMessage); err != nil {
		return err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	peer, ok := n.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
	}

	msg := NewMessage(MessageTypeStatus, buf.Bytes())

	return peer.Send(msg.Bytes())
}

func (n *Node) processBlock(b *types.Block) error {
	if err := n.chain.AddBlock(b); err != nil {
		n.Logger.Log("error", err.Error())
		return err
	}

	go n.broadcastBlock(b)

	return nil
}

func (n *Node) processTransaction(tx *types.Transaction) error {
	hash := tx.GetHash(types.TxHasher{})

	if n.mempool.Contains(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	// s.Logger.Log(
	// 	"msg", "adding new tx to mempool",
	// 	"hash", hash,
	// 	"mempoolPending", s.mempool.PendingCount(),
	// )

	go n.broadcastTx(tx)

	n.mempool.Add(tx)

	return nil
}

// 네트워크에서 가장 높은 블록 높이에 있을 때 계속 동기화되지 않도록 하는 방법을 찾아야 함.
func (n *Node) requestBlocksLoop(peer net.Addr) error {
	ticker := time.NewTicker(3 * time.Second)

	for {
		ourHeight := n.chain.Height()

		n.Logger.Log("msg", "👋 requesting block height from", ourHeight+1)

		getBlocksMessage := &GetBlocksMessage{
			From: ourHeight + 1,
			To:   0,
		}

		buf := new(bytes.Buffer)
		if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err != nil {
			return err
		}

		n.mu.RLock()
		defer n.mu.RUnlock()

		msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
		peer, ok := n.peerMap[peer]
		if !ok {
			return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
		}

		if err := peer.Send(msg.Bytes()); err != nil {
			n.Logger.Log("error", "failed to send to peer", "err", err, "peer", peer)
		}

		<-ticker.C
	}
}

func (n *Node) broadcastBlock(b *types.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(types.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return n.broadcast(msg.Bytes())
}

func (n *Node) broadcastTx(tx *types.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(types.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return n.broadcast(msg.Bytes())
}

func (n *Node) createNewBlock() error {
	currentHeader, err := n.chain.GetHeader(n.chain.Height())
	if err != nil {
		return err
	}

	// 우선은 멤풀에 있는 모든 트랜잭션을 블록에 담고 추후 수정 예정.
	// 트랜잭션을 아직 구체화하지 않았기 때문.
	txx := n.mempool.Pending()

	block, err := types.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*n.PrivateKey); err != nil {
		return err
	}

	if err := n.chain.AddBlock(block); err != nil {
		return err
	}

	n.mempool.ClearPending()

	go n.broadcastBlock(block)

	return nil
}

func genesisBlock() *types.Block {
	header := &types.Header{
		Version:   1,
		DataHash:  common.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := types.NewBlock(header, nil)

	coinbase := crypto.PublicKey{}
	tx := types.NewTransaction(nil)
	tx.From = coinbase
	tx.To = coinbase
	tx.Value = 10_000_000
	b.Transactions = append(b.Transactions, tx)

	privKey := crypto.GeneratePrivateKey()
	if err := b.Sign(privKey); err != nil {
		panic(err)
	}

	return b
}
