package core

import (
	"github.com/barreleye-labs/barreleye/barreldb/query"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) GetBlockFromDB(hash common.Hash, number uint64) *types.Block {
	query.GetBlock()
	return nil
}

func (bc *Blockchain) GetLastBlockFromDB(hash common.Hash, number uint64) *types.Block {
	query.GetLastBlock()
	return nil
}
