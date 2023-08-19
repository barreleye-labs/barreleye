package core

import (
	"crypto/sha256"

	"github.com/barreleye-labs/barreleye/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}
