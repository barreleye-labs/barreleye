package common

import "errors"

var (
	LevelDBNotFoundError         = "leveldb: not found"
	ErrBlockKnown                = errors.New("block already known")
	ErrTransactionAlreadyPending = errors.New("this transaction is already pending transaction")
	ErrBlockTooHigh              = errors.New("block Too high")
)
