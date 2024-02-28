package node

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/barreleye-labs/barreleye/core/types"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/barreleye-labs/barreleye/core"
	"github.com/barreleye-labs/barreleye/restful"
	"github.com/go-kit/log"
)

var defaultBlockTime = 7 * time.Second

type NodeOpts struct {
	APIListenAddr string
	SeedNodes     []string
	ListenAddr    string
	TCPTransport  *TCPTransport
	Name          string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	BlockTime     time.Duration
	PrivateKey    *types.PrivateKey
}

type Node struct {
	TCPTransport *TCPTransport
	peerCh       chan *TCPPeer

	mu      sync.RWMutex
	peerMap map[net.Addr]*TCPPeer

	NodeOpts
	txPool       *TxPool
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
		opts.RPCDecodeFunc = DecodeRPCDefaultFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "🕰", log.DefaultTimestampUTC)
	}

	chain, err := core.NewBlockchain(opts.Logger, opts.PrivateKey)
	if err != nil {
		return nil, err
	}

	txChan := make(chan *types.Transaction)

	if len(opts.APIListenAddr) > 0 {
		apiNodeCfg := restful.ServerConfig{
			Logger:     opts.Logger,
			ListenAddr: opts.APIListenAddr,
		}
		apiNode := restful.NewServer(apiNodeCfg, chain, txChan, opts.PrivateKey)
		go apiNode.Start()

		_ = opts.Logger.Log("msg", "HTTP API server running", "port", opts.APIListenAddr)
	}

	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddr, peerCh)

	s := &Node{
		TCPTransport: tr,
		peerCh:       peerCh,
		peerMap:      make(map[net.Addr]*TCPPeer),
		NodeOpts:     opts,
		chain:        chain,
		txPool:       NewTxPool(1000),
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

	_ = n.Logger.Log("msg", "🤝 Ready to connect with peers", "port", n.ListenAddr, "name", n.Name)

free:
	for {
		select {
		case peer := <-n.peerCh:
			n.peerMap[peer.conn.RemoteAddr()] = peer

			go peer.readLoop(n.rpcCh)

			if err := n.sendChainInfoRequestMessage(peer); err != nil {
				_ = n.Logger.Log("err", err)
				continue
			}

			_ = n.Logger.Log("msg", "🙋 connected peer", "peer", peer.conn.RemoteAddr())

		case tx := <-n.txChan:
			if err := n.handleTransaction(tx); err != nil {
				_ = n.Logger.Log("process TX error", err)
			}

		case rpc := <-n.rpcCh:
			msg, err := n.RPCDecodeFunc(rpc)
			if err != nil {
				_ = n.Logger.Log("RPC error", err)
				continue
			}

			if err = n.RPCProcessor.HandleMessage(msg); err != nil {
				if !errors.Is(err, core.ErrBlockKnown) && !errors.Is(err, core.ErrTransactionAlreadyPending) {
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

func (n *Node) HandleMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *types.Transaction:
		return n.handleTransaction(t)
	case *types.Block:
		return n.handleBlock(t)
	case *ChainInfoRequestMessage:
		return n.handleChainInfoRequestMessage(msg.From)
	case *ChainInfoResponseMessage:
		return n.handleChainInfoResponseMessage(msg.From, t)
	case *BlockRequestMessage:
		return n.handleBlockRequestMessage(msg.From, t)
	case *BlockResponseMessage:
		return n.handleBlockResponseMessage(msg.From, t)
	}

	return nil
}

func (n *Node) handleBlock(b *types.Block) error {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	n.miningTicker.Reset(n.BlockTime + time.Duration(r.Intn(7))*time.Second)
	if err := n.chain.LinkBlock(b); err != nil {
		//_ = n.Logger.Log("error", err.Error())
		return err
	}

	go n.broadcastBlock(b)

	return nil
}

func (n *Node) handleTransaction(tx *types.Transaction) error {
	if err := tx.Verify(); err != nil {
		return err
	}

	if err := n.txPool.Add(tx, n.chain); err != nil {
		return err
	}

	go n.broadcastTx(tx)

	hash := tx.GetHash()

	_ = n.Logger.Log(
		"msg", "🗃️ adding new tx to txpool",
		"hash", hash,
		"PendingCount", n.txPool.PendingCount(),
	)

	return nil
}

func (n *Node) handleBlockRequestMessage(from net.Addr, data *BlockRequestMessage) error {
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

	_ = n.Logger.Log("msg", "✉️ send chain info request message", "to", peer.conn.RemoteAddr())
	return nil
}

func (n *Node) broadcast(payload []byte) error {
	n.mu.RLock()
	defer n.mu.RUnlock()
	for netAddr, peer := range n.peerMap {
		if err := peer.Send(payload); err != nil {
			if err = peer.Close(); err != nil {
				return err
			}
			delete(n.peerMap, netAddr)
			fmt.Printf("communication with peer has been lost and no further transmissions will be made..\nwrite error: ", err)
		}
	}
	return nil
}

func (n *Node) handleBlockResponseMessage(from net.Addr, data *BlockResponseMessage) error {
	_ = n.Logger.Log("msg", "📦 received the requested block", "height:", data.Block.Height, "from", from)

	if data.Block == nil {
		return fmt.Errorf("no block in block response message")
	}

	if err := n.chain.LinkBlock(data.Block); err != nil {
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

func (n *Node) handleChainInfoResponseMessage(from net.Addr, data *ChainInfoResponseMessage) error {
	_ = n.Logger.Log("msg", "📬 received chain info response message", "from", from, "height", data.CurrentHeight)

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

func (n *Node) handleChainInfoRequestMessage(from net.Addr) error {
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
		To:            n.Name,
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

	if err := peer.Send(msg.Bytes()); err != nil {
		return err
	}

	_ = n.Logger.Log("msg", "✉️ send chain info response message", "to", peer.conn.RemoteAddr())
	return nil
}

// 네트워크에서 가장 높은 블록 높이에 있을 때 계속 동기화되지 않도록 하는 방법을 찾아야 함.
func (n *Node) sendBlockRequestMessage(peerAddr net.Addr, blockNumber int32) error {
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

	_ = n.Logger.Log("msg", "✉️ send block request message", "height", blockNumber)
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

	txs := n.txPool.Pending()

	for i := 0; i < len(txs); i++ {
		txProcessed, err := n.chain.ReadTxByHash(txs[i].Hash)
		if err != nil {
			return err
		}

		if txProcessed != nil {
			txs[i] = txs[len(txs)-1]
			txs = txs[:len(txs)-1]
			i--
			continue
		}

		accountNonce, err := n.chain.ReadAccountNonceByAddress(txs[i].From)
		if err != nil {
			return err
		}

		nonce := uint64(0)
		if accountNonce != nil {
			nonce = *accountNonce
		}

		if nonce != txs[i].Nonce {
			txs[i] = txs[len(txs)-1]
			txs = txs[:len(txs)-1]
			i--
		}
	}

	block, err := types.NewBlockFromPrevHeader(lastHeader, txs)
	if err != nil {
		return err
	}

	if err := block.Sign(*n.PrivateKey); err != nil {
		return err
	}

	_ = n.Logger.Log("msg", "🍀 block mining success")

	if err := n.chain.LinkBlock(block); err != nil {
		return err
	}

	n.txPool.ClearPending()

	go n.broadcastBlock(block)

	return nil
}
