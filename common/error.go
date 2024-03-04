package common

import "errors"

var (
	LevelDBNotFoundError         = "leveldb: not found"
	ErrBlockKnown                = errors.New("block already known")
	ErrTransactionAlreadyPending = errors.New("this transaction is already pending transaction")
	ErrBlockTooHigh              = errors.New("block Too high")
	ErrPrevBlockMismatch         = errors.New("previous block hash of the block to be connected does not match the current block hash")
)
