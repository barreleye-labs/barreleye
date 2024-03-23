package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/common/util"
	"log"
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
	nonce := util.Uint64ToBytes(tx.Nonce)
	from := tx.From.ToSlice()
	to := tx.To.ToSlice()
	value := util.Uint64ToBytes(tx.Value)
	data := tx.Data
	buf := new(bytes.Buffer)
	_ = binary.Write(buf, binary.LittleEndian, nonce)
	_ = binary.Write(buf, binary.LittleEndian, from)
	_ = binary.Write(buf, binary.LittleEndian, to)
	_ = binary.Write(buf, binary.LittleEndian, value)
	_ = binary.Write(buf, binary.LittleEndian, data)

	msgHash := fmt.Sprintf(
		"%x",
		sha256.Sum256([]byte(hex.EncodeToString(buf.Bytes()))),
	)

	message, hashDecodeError := hex.DecodeString(msgHash)
	if hashDecodeError != nil {
		log.Println(hashDecodeError)
		panic("internal server error")
	}

	return common.HashFromBytes(message)
}
