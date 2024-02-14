package core

import (
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

	lastBlockHeight, err := v.bc.ReadLastBlockHeight()

	if *lastBlockHeight >= b.Height {
		return ErrBlockKnown
	}

	if b.Height != *lastBlockHeight+1 {
		return fmt.Errorf("block (%s) with height (%d) is too high => current height (%d)", b.GetHash(), b.Height, *lastBlockHeight)
	}

	prevHeader, err := v.bc.ReadHeaderByHeight(b.Height - 1)
	if err != nil {
		return err
	}

	hash := types.BlockHasher{}.Hash(prevHeader)

	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	if err := b.Verify(); err != nil {
		return err
	}

	return nil
}
