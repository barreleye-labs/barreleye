package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"time"
)

type Header struct {
	Version       uint32
	DataHash      common.Hash
	PrevBlockHash common.Hash
	Height        int32
	Timestamp     int64
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	enc.Encode(h)

	return buf.Bytes()
}

type Block struct {
	*Header

	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// Cached version of the header hash
	Hash common.Hash
}

func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	block := &Block{
		Header:       h,
		Transactions: txx,
	}
	block.Hash = block.GetHash()
	return block, nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, txx)
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
	hash, _ := CalculateDataHash(b.Transactions)
	b.DataHash = hash
}

func (b *Block) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(b.GetHash().ToSlice())
	if err != nil {
		return err
	}

	b.Validator = privKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.GetHash().ToSlice()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	dataHash, err := CalculateDataHash(b.Transactions)
	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.GetHash())
	}

	return nil
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	return dec.Decode(b)
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	return enc.Encode(b)
}

func (b *Block) GetHash() common.Hash {
	if b.Hash.IsZero() {
		b.Hash = BlockHasher{}.Hash(b.Header)
	}

	return b.Hash
}

func CalculateDataHash(txx []*Transaction) (hash common.Hash, err error) {
	buf := []byte{}
	for _, tx := range txx {
		buf = append(buf, tx.GetHash().ToSlice()...)
	}
	hash = sha256.Sum256(buf)

	return
}

func init() {
	gob.Register(secp256k1.S256())
}
