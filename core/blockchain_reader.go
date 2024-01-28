package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) GetBlockByHashFromDB(hash common.Hash) (*types.Block, error) {
	block, err := bc.db.GetBlock(hash)
	if err != nil {
		return nil, fmt.Errorf("fail to get block")
	}

	return block, nil
}
