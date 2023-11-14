package core

import (
	"crypto/sha256"
	"encoding/binary"

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

type TxHasher struct{}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, uint64(tx.Nounce))
	data := append(buf, tx.Data...)
	
	return types.Hash(sha256.Sum256(data))
}
