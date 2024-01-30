package core

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) WriteBlockWithHash(hash common.Hash, block *types.Block) error {
	if err := bc.db.InsertBlockWithHash(hash, block); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteBlockWithHeight(height uint32, block *types.Block) error {
	if err := bc.db.InsertBlockWithHeight(height, block); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteLastBlock(block *types.Block) error {
	if err := bc.db.InsertLastBlock(block); err != nil {
		return err
	}
	return nil
}
