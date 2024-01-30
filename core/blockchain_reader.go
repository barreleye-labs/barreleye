package core

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) ReadBlockByHash(hash common.Hash) (*types.Block, error) {
	block, err := bc.db.SelectBlockByHash(hash)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (bc *Blockchain) ReadBlockByHeight(height uint32) (*types.Block, error) {
	block, err := bc.db.SelectBlockByHeight(height)
	if err != nil {
		return nil, err
	}

	return block, nil
}

func (bc *Blockchain) ReadBlocks(page int, size int) ([]*types.Block, error) {

	offset := (page - 1) * size

	lastBlock, err := bc.ReadLastBlock()
	if err != nil {
		return nil, err
	}

	blocks := []*types.Block{}

	lastBlockHeight := int(lastBlock.Height)
	start := lastBlockHeight - offset
	if start < 0 {
		return blocks, nil
	}

	end := start - size
	if end < -1 {
		end = -1
	}

	for i := start; i > end; i-- {
		block, err := bc.ReadBlockByHeight(uint32(i))
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (bc *Blockchain) ReadBlocksByHash(hash common.Hash, size int) ([]*types.Block, error) {
	blocks := []*types.Block{}
	for i := 0; i < size; i++ {
		block, _ := bc.ReadBlockByHash(hash)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
		hash = block.PrevBlockHash
	}
	return blocks, nil
}

func (bc *Blockchain) ReadLastBlock() (*types.Block, error) {
	block, err := bc.db.SelectLastBlock()
	if err != nil {
		return nil, err
	}

	return block, nil
}
