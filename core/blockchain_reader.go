package core

import (
	"github.com/barreleye-labs/barreleye/barreldb/query"
	"github.com/barreleye-labs/barreleye/common"
)

func (bc *Blockchain) GetBlockFromDB(hash common.Hash, number uint64) *Block {
	query.GetBlock()
	return nil
}

func (bc *Blockchain) GetLastBlockFromDB(hash common.Hash, number uint64) *Block {
	query.GetLastBlock()
	return nil
}
