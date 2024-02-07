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

	_ = binary.Write(buf, binary.LittleEndian, header.Version)
	_ = binary.Write(buf, binary.LittleEndian, header.DataHash)
	_ = binary.Write(buf, binary.LittleEndian, header.PrevBlockHash)
	_ = binary.Write(buf, binary.LittleEndian, header.Height)
	_ = binary.Write(buf, binary.LittleEndian, header.Timestamp)

	return sha256.Sum256(buf.Bytes())
}

type TxHasher struct{}

// Hash will hash the whole bytes of the TX no exception.
func (TxHasher) Hash(tx *Transaction) common.Hash {
	buf := new(bytes.Buffer)

	_ = binary.Write(buf, binary.LittleEndian, tx.Nonce)
	_ = binary.Write(buf, binary.LittleEndian, tx.From)
	_ = binary.Write(buf, binary.LittleEndian, tx.To)
	_ = binary.Write(buf, binary.LittleEndian, tx.Value)
	_ = binary.Write(buf, binary.LittleEndian, tx.Data)
	_ = binary.Write(buf, binary.LittleEndian, tx.Signer)

	return sha256.Sum256(buf.Bytes())
}
