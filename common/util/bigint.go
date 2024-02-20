package util

import "math/big"

func Uint64ToBytes(number uint64) []byte {
	bigint := new(big.Int)
	bigint.SetUint64(number)
	bytes := bigint.Bytes()
	if len(bytes) == 0 {
		bytes = []byte{0}
	}
	return bytes
}

func Int64ToBytes(number int64) []byte {
	bigint := new(big.Int)
	bigint.SetInt64(number)
	return bigint.Bytes()
}
