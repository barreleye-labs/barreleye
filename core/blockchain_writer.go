package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) CreateBlock(hash common.Hash, block *types.Block) error {
	if err := bc.db.CreateBlock(hash, block); err != nil {
		return fmt.Errorf("fail to create block")
	}
	return nil
}
