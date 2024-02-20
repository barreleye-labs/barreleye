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
	MessageTypeTx            MessageType = 0x1
	MessageTypeBlock         MessageType = 0x2
	MessageTypeBlockRequest  MessageType = 0x3
	MessageTypeStatus        MessageType = 0x4
	MessageTypeGetStatus     MessageType = 0x5
	MessageTypeBlockResponse MessageType = 0x6
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
	gob.NewEncoder(buf).Encode(msg)
	return buf.Bytes()
}

type DecodedMessage struct {
	From net.Addr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := Message{}
	if err := gob.NewDecoder(rpc.Payload).Decode(&msg); err != nil {
		return nil, fmt.Errorf("failed to decode message from %s: %s", rpc.From, err)
	}

	// fmt.Printf("receiving message: %+v\n", msg)

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

	case MessageTypeGetStatus:
		return &DecodedMessage{
			From: rpc.From,
			Data: &ChainInfoRequestMessage{},
		}, nil

	case MessageTypeStatus:
		statusMessage := new(ChainInfoResponseMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMessage); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: statusMessage,
		}, nil

	case MessageTypeBlockRequest:
		getBlocks := new(BlockRequestMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocks); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: getBlocks,
		}, nil

	case MessageTypeBlockResponse:
		blocks := new(BlockResponseMessage)
		if err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blocks); err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blocks,
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
