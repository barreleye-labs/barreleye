package common

import (
	"encoding/hex"
	"fmt"
)

const (
	AddressLength = 20
)

type Address [AddressLength]byte

func (a Address) ToSlice() []byte {
	b := make([]byte, 20)
	for i := 0; i < 20; i++ {
		b[i] = a[i]
	}
	return b
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func (a Address) Equal(address Address) bool {
	if a.String() == address.String() {
		return true
	}
	return false
}

func NewAddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 32", len(b))
		panic(msg)
	}

	var value [20]uint8
	for i := 0; i < 20; i++ {
		value[i] = b[i]
	}

	return value
}
