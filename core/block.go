package core

import (
	"io"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/barreleye-labs/barreleye/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Height        uint32
	Timestamp     uint64
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// Cached version of the header hash
	hash types.Hash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

func (b *Block) Decode(r io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(w, b)
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}

// func (b *Block) Hash() types.Hash {
// 	buf := &bytes.Buffer{}
// 	b.Header.EncodeBinary(buf)

// 	b.hash = types.Hash(sha256.Sum256(buf.Bytes()))

// 	return b.hash
// }
