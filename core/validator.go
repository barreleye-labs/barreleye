package core

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/barreleye-labs/barreleye/core/types"
)

var ErrBlockKnown = errors.New("block already known")

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
		return ErrBlockKnown
	}

	if err = b.Verify(); err != nil {
		return err
	}

	// 블록 높이가 같은데 해시가 다를 경우 기존 블록을 버리고 받은 블록으로 덮어 씌움.
	if lastBlock.Height == b.Height {
		if lastBlock.Hash.String() != b.Hash.String() && bytes.Compare(lastBlock.Hash.ToSlice(), b.Hash.ToSlice()) == 1 {
			_ = v.bc.logger.Log("msg", "block replacement")
			if err = v.bc.RemoveLastBlock(); err != nil {
				return err
			}
			return nil
		}
		return ErrBlockKnown
	}

	if lastBlock.Height+1 != b.Height {
		return fmt.Errorf("block (%s) with height (%d) is too high => current height (%d)", b.GetHash(), b.Height, lastBlock.Height)
	}

	prevHeader, err := v.bc.ReadHeaderByHeight(b.Height - 1)
	if err != nil {
		return err
	}

	hash := types.BlockHasher{}.Hash(prevHeader)

	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	return nil
}
