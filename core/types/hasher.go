package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"github.com/barreleye-labs/barreleye/common"
)

type Hasher[T any] interface {
	Hash(T) common.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(header *Header) common.Hash {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, header.Version)
	binary.Write(buf, binary.LittleEndian, header.DataHash)
	binary.Write(buf, binary.LittleEndian, header.PrevBlockHash)
	binary.Write(buf, binary.LittleEndian, header.Height)
	binary.Write(buf, binary.LittleEndian, header.Timestamp)

	return sha256.Sum256(buf.Bytes())
}

type TxHasher struct{}

// Hash will hash the whole bytes of the TX no exception.
func (TxHasher) Hash(tx *Transaction) common.Hash {
	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, tx.Data)
	binary.Write(buf, binary.LittleEndian, tx.To)
	binary.Write(buf, binary.LittleEndian, tx.Value)
	binary.Write(buf, binary.LittleEndian, tx.From)
	binary.Write(buf, binary.LittleEndian, tx.Nonce)

	return sha256.Sum256(buf.Bytes())
}
