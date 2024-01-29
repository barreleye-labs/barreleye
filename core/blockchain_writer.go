package core

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) CreateBlockWithHash(hash common.Hash, block *types.Block) error {
	if err := bc.db.CreateBlockWithHash(hash, block); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) CreateBlockWithHeight(height uint32, block *types.Block) error {
	if err := bc.db.CreateBlockWithHeight(height, block); err != nil {
		return err
	}
	return nil
}
