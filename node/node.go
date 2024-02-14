package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/barreleye-labs/barreleye/core/types"
	"math/rand"
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
	mempool      *TxPool
	chain        *core.Blockchain
	isValidator  bool
	rpcCh        chan RPC
	quitCh       chan struct{}
	txChan       chan *types.Transaction
	miningTicker *time.Ticker

	peersBlockHeightUntilSync int32
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
		opts.Logger = log.With(opts.Logger, "🕰", log.DefaultTimestampUTC)
	}

	var genesis *types.Block = nil
	if opts.ID == "GENESIS-NODE" {
		genesis = CreateGenesisBlock(opts.PrivateKey)
		_ = opts.Logger.Log("msg", "🌞 create genesis block")
	}

	chain, err := core.NewBlockchain(opts.Logger, opts.PrivateKey, genesis)
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
		miningTicker: time.NewTicker(opts.BlockTime),
	}

	s.TCPTransport.peerCh = peerCh

	if s.RPCProcessor == nil {
		s.RPCProcessor = s
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

	_ = n.Logger.Log("msg", "🤝 accepting TCP connection on", "addr", n.ListenAddr, "id", n.ID)

free:
	for {
		select {
		case peer := <-n.peerCh:
			n.peerMap[peer.conn.RemoteAddr()] = peer

			go peer.readLoop(n.rpcCh)

			if err := n.sendChainInfoRequestMessage(peer); err != nil {
				n.Logger.Log("err", err)
				continue
			}

			_ = n.Logger.Log("msg", "🙋 peer added to the Node", "outgoing", peer.Outgoing, "addr", peer.conn.RemoteAddr())

		case tx := <-n.txChan:
			if err := n.processTransaction(tx); err != nil {
				_ = n.Logger.Log("process TX error", err)
			}

		case rpc := <-n.rpcCh:
			msg, err := n.RPCDecodeFunc(rpc)
			if err != nil {
				_ = n.Logger.Log("RPC error", err)
				continue
			}

			if err := n.RPCProcessor.ProcessMessage(msg); err != nil {
				if err != core.ErrBlockKnown {
					_ = n.Logger.Log("error", err)
				}
			}

		case <-n.quitCh:
			break free
		}
	}

	_ = n.Logger.Log("msg", "Node is shutting down")
}

func (n *Node) mine() error {
	_ = n.Logger.Log("msg", "start mining using POR(proof of random)", "blockTime", n.BlockTime)

	for {
		//height, err := n.chain.ReadLastBlockHeight()
		//if err != nil {
		//	return err
		//}
		//
		//if n.peersBlockHeightUntilSync > *height {
		//	continue
		//}

		<-n.miningTicker.C

		if err := n.sealBlock(); err != nil {
			_ = n.Logger.Log("sealing block error", err)
		}
	}
}

func (n *Node) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *types.Transaction:
		return n.processTransaction(t)
	case *types.Block:
		return n.processBlock(t)
	case *ChainInfoRequestMessage:
		return n.processChainInfoRequestMessage(msg.From)
	case *ChainInfoResponseMessage:
		return n.processChainInfoResponseMessage(msg.From, t)
	case *BlockRequestMessage:
		return n.processBlockRequestMessage(msg.From, t)
	case *BlockResponseMessage:
		return n.processBlockResponseMessage(msg.From, t)
	}

	return nil
}

func (n *Node) processBlock(b *types.Block) error {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	n.miningTicker.Reset(n.BlockTime + time.Duration(r.Intn(2))*time.Second)
	if err := n.chain.AddBlock(b); err != nil {
		n.Logger.Log("error", err.Error())
		return err
	}

	go n.broadcastBlock(b)

	return nil
}

func (n *Node) processTransaction(tx *types.Transaction) error {
	hash := tx.GetHash()

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

func (n *Node) processBlockRequestMessage(from net.Addr, data *BlockRequestMessage) error {
	_ = n.Logger.Log("msg", "📬 received blockRequest message", "from", from)

	height, err := n.chain.ReadLastBlockHeight()
	if err != nil {
		return err
	}

	if *height < data.Height {
		return fmt.Errorf("requested block number %d is higher compared to block number %d in this chain", data.Height, height)
	}

	block, err := n.chain.ReadBlockByHeight(data.Height)
	if err != nil {
		return err
	}

	if block == nil {
		return fmt.Errorf("not found block")
	}

	blockResponseMsg := &BlockResponseMessage{
		Block: block,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(blockResponseMsg); err != nil {
		return err
	}

	n.mu.RLock()
	defer n.mu.RUnlock()

	msg := NewMessage(MessageTypeBlockResponse, buf.Bytes())
	peer, ok := n.peerMap[from]
	if !ok {
		return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
	}

	return peer.Send(msg.Bytes())
}

func (n *Node) sendChainInfoRequestMessage(peer *TCPPeer) error {
	var (
		getStatusMsg = new(ChainInfoRequestMessage)
		buf          = new(bytes.Buffer)
	)

	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())

	if err := peer.Send(msg.Bytes()); err != nil {
		return err
	}

	_ = n.Logger.Log("msg", "👋 requesting chain info request message", "to", peer.conn.RemoteAddr())
	return nil
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

func (n *Node) processBlockResponseMessage(from net.Addr, data *BlockResponseMessage) error {
	_ = n.Logger.Log("msg", "📦 received the requested block", "height:", data.Block.Height, "from", from)

	if data.Block == nil {
		return fmt.Errorf("no block in block response message")
	}

	if err := n.chain.AddBlock(data.Block); err != nil {
		_ = n.Logger.Log("error", err.Error())
		return err
	}

	if n.peersBlockHeightUntilSync > data.Block.Height {
		if err := n.sendBlockRequestMessage(from, data.Block.Height+1); err != nil {
			return err
		}
	} else if n.peersBlockHeightUntilSync == data.Block.Height {
		peer := n.peerMap[from]
		if err := n.sendChainInfoRequestMessage(peer); err != nil {
			return err
		}
	}

	return nil
}

func (n *Node) processChainInfoResponseMessage(from net.Addr, data *ChainInfoResponseMessage) error {
	n.Logger.Log("msg", "📬 received chain info response message", "from", from)

	height, err := n.chain.ReadLastBlockHeight()
	if err != nil {
		return err
	}

	// 전달 받은 블록 높이보다 현재 나의 블록체인의 블록 높이가 같거나 클 경우.
	if data.CurrentHeight <= *height {
		n.Logger.Log("msg", "already sync", "this node height", height, "network height", data.CurrentHeight, "addr", from)
		go n.mine()
		return nil
	}

	n.peersBlockHeightUntilSync = data.CurrentHeight

	if err = n.sendBlockRequestMessage(from, *height+1); err != nil {
		return err
	}

	return nil
}

func (n *Node) processChainInfoRequestMessage(from net.Addr) error {
	_ = n.Logger.Log("msg", "📬 received chain info request message", "from", from)

	height, err := n.chain.ReadLastBlockHeight()
	if err != nil {
		return err
	}
	return n.sendChainInfoResponseMessage(from, *height)
}

func (n *Node) sendChainInfoResponseMessage(from net.Addr, height int32) error {
	chainInfoResponseMessage := &ChainInfoResponseMessage{
		CurrentHeight: height,
		ID:            n.ID,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(chainInfoResponseMessage); err != nil {
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

// 네트워크에서 가장 높은 블록 높이에 있을 때 계속 동기화되지 않도록 하는 방법을 찾아야 함.
func (n *Node) sendBlockRequestMessage(peerAddr net.Addr, blockNumber int32) error {
	_ = n.Logger.Log("msg", "👋 requesting block height from", blockNumber)

	blockRequestMessage := &BlockRequestMessage{
		Height: blockNumber,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(blockRequestMessage); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlockRequest, buf.Bytes())
	peer, ok := n.peerMap[peerAddr]
	if !ok {
		return fmt.Errorf("peer %s not known", peer.conn.RemoteAddr())
	}

	if err := peer.Send(msg.Bytes()); err != nil {
		_ = n.Logger.Log("error", "failed to send to peer", "err", err, "peer", peer)
	}

	return nil
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

func (n *Node) sealBlock() error {
	lastHeader, err := n.chain.ReadLastHeader()
	if err != nil {
		return err
	}

	if lastHeader == nil {
		return fmt.Errorf("can not seal the block without genesis block")
	}

	// 우선은 멤풀에 있는 모든 트랜잭션을 블록에 담고 추후 수정 예정.
	// 트랜잭션을 아직 구체화하지 않았기 때문.
	txx := n.mempool.Pending()

	block, err := types.NewBlockFromPrevHeader(lastHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*n.PrivateKey); err != nil {
		return err
	}

	_ = n.Logger.Log("msg", "🍀 block mining success")

	if err := n.chain.AddBlock(block); err != nil {
		return err
	}

	n.mempool.ClearPending()

	go n.broadcastBlock(block)

	return nil
}

func CreateGenesisBlock(privateKey *crypto.PrivateKey) *types.Block {
	coinbase := privateKey.PublicKey()

	tx := &types.Transaction{
		Nonce: 171, //ab
		From:  coinbase.Address(),
		To:    coinbase.Address(),
		Value: 171, //ab
		Data:  []byte{171},
	}

	if err := tx.Sign(*privateKey); err != nil {
		panic(err)
	}

	header := &types.Header{
		Version:   1,
		Height:    0,
		Timestamp: time.Now().UnixNano(),
	}

	b, _ := types.NewBlock(header, nil)

	b.Transactions = append(b.Transactions, tx)
	hash, _ := types.CalculateDataHash(b.Transactions)
	b.DataHash = hash

	if err := b.Sign(*privateKey); err != nil {
		panic(err)
	}
	return b
}
