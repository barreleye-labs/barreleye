package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

type Validator interface {
	ValidateBlock(*types.Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidateBlock(b *types.Block) error {
	if b.Height == 0 {
		return nil
	}

	lastBlock, err := v.bc.ReadLastBlock()
	if err != nil {
		return err
	}

	if lastBlock == nil {
		return fmt.Errorf("lastBlock is nil")
	}

	if lastBlock.Height > b.Height {
		return common.ErrBlockKnown
	}

	if lastBlock.Height+1 < b.Height {
		return common.ErrBlockTooHigh
	}

	prevHeader, err := v.bc.ReadHeaderByHeight(b.Height - 1)
	if err != nil {
		return err
	}

	hash := types.BlockHasher{}.Hash(prevHeader)

	if hash != b.PrevBlockHash {
		return common.ErrPrevBlockMismatch
	}

	if err = b.Verify(); err != nil {
		return err
	}

	// 블록 높이가 같은 다른 블록을 수신한 경우 해시값이 작은 블록을 선택함.
	if lastBlock.Height == b.Height {
		if !lastBlock.Hash.Equal(b.Hash) && lastBlock.Hash.Compare(b.Hash) == 1 {
			_ = v.bc.logger.Log("msg", "block replacement")
			if err = v.bc.RemoveLastBlock(); err != nil {
				return err
			}
			return nil
		}
		return common.ErrBlockKnown
	}
	return nil
}
