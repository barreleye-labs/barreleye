package node

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"io"
	"net"

	"github.com/sirupsen/logrus"
)

type MessageType byte

const (
	MessageTypeTx                MessageType = 0x1
	MessageTypeBlock             MessageType = 0x2
	MessageTypeChainInfoResponse MessageType = 0x3
	MessageTypeChainInfoRequest  MessageType = 0x4
	MessageTypeBlockRequest      MessageType = 0x5
	MessageTypeBlockResponse     MessageType = 0x6
	MessageTypeBlockHashRequest  MessageType = 0x7
	MessageTypeBlockHashResponse MessageType = 0x8
)

type RPC struct {
	From    net.Addr //string
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func (msg *Message) Bytes() []byte {
	buf := &bytes.Buffer{}
	_ = gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type DecodedMessage struct {
	From net.Addr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DecodeRPCDefaultFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	logrus.WithFields(logrus.Fields{
		"from": rpc.From,
		"type": msg.Header,
	}).Debug("new incoming message")

	switch msg.Header {
	case MessageTypeTx:
		tx := new(types.Transaction)
		if err := tx.Decode(types.NewGobTxDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil

	case MessageTypeBlock:
		block := new(types.Block)
		if err := block.Decode(types.NewGobBlockDecoder(bytes.NewReader(msg.Data))); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: block,
		}, nil

	case MessageTypeChainInfoRequest:
		return &DecodedMessage{
			From: rpc.From,
			Data: &ChainInfoRequestMessage{},
		}, nil

	case MessageTypeChainInfoResponse:
		chainInfoResponseMessage := new(ChainInfoResponseMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(chainInfoResponseMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: chainInfoResponseMessage,
		}, nil

	case MessageTypeBlockRequest:
		blockRequestMessage := new(BlockRequestMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blockRequestMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blockRequestMessage,
		}, nil

	case MessageTypeBlockResponse:
		blockResponseMessage := new(BlockResponseMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blockResponseMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blockResponseMessage,
		}, nil

	case MessageTypeBlockHashRequest:
		blockHashRequestMessage := new(BlockHashRequestMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blockHashRequestMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blockHashRequestMessage,
		}, nil

	case MessageTypeBlockHashResponse:
		blockHashResponseMessage := new(BlockHashResponseMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blockHashResponseMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blockHashResponseMessage,
		}, nil
	default:
		return nil, fmt.Errorf("invalid message header %x", msg.Header)
	}

}

type RPCProcessor interface {
	HandleMessage(*DecodedMessage) error
}

func init() {
	gob.Register(secp256k1.S256())
}
