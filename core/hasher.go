package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/barreleye-labs/barreleye/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Block) types.Hash {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(b.Header); err != nil {
		panic(err)
	}

	h := sha256.Sum256(buf.Bytes())
	return types.Hash(h)
}
