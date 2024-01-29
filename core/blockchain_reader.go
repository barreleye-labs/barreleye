package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) GetBlockByHashFromDB(hash common.Hash) (*types.Block, error) {
	block, err := bc.db.GetBlock(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get block")
	}

	return block, nil
}

func (bc *Blockchain) GetBlocks(hash common.Hash, size int) ([]*types.Block, error) {
	blocks := []*types.Block{}
	for i := 0; i < size; i++ {
		block, _ := bc.GetBlockByHashFromDB(hash)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
		hash = block.PrevBlockHash
	}
	return blocks, nil
}
