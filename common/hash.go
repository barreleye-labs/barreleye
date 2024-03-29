package common

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

const (
	HashLength = 32
)

type Hash [HashLength]byte

func (h Hash) IsZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] != 0 {
			return false
		}
	}
	return true
}

func (h Hash) ToSlice() []byte {
	b := make([]byte, 32)
	for i := 0; i < 32; i++ {
		b[i] = h[i]
	}
	return b
}

func (h Hash) Equal(hash Hash) bool {
	return bytes.Equal(h.ToSlice(), hash.ToSlice())
}

func (h Hash) Compare(hash Hash) int {
	return bytes.Compare(h.ToSlice(), hash.ToSlice())
}

func (h Hash) String() string {
	return hex.EncodeToString(h.ToSlice())
}

func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}
	var value [32]uint8
	for i := 0; i < 32; i++ {
		value[i] = b[i]
	}

	return value
}
