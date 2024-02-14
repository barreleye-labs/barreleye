package core

import (
	"fmt"
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

func (bc *Blockchain) ReadBlockByHeight(height int32) (*types.Block, error) {
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

	if lastBlock == nil {
		return blocks, nil
	}

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
		block, err := bc.ReadBlockByHeight(int32(i))
		if err != nil {
			return nil, err
		}
		if block == nil {
			return nil, fmt.Errorf("block %d is nil", i)
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

// 사용시 수정 필요함.
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

func (bc *Blockchain) ReadLastBlockHeight() (*int32, error) {
	block, err := bc.ReadLastBlock()
	if err != nil {
		return nil, err
	}

	if block == nil {
		height := int32(-1)
		return &height, nil
	}

	return &block.Height, nil
}

func (bc *Blockchain) ReadHeaderByHash(hash common.Hash) (*types.Header, error) {
	header, err := bc.db.SelectHeaderByHash(hash)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (bc *Blockchain) ReadHeaderByHeight(height int32) (*types.Header, error) {
	header, err := bc.db.SelectHeaderByHeight(height)
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (bc *Blockchain) ReadHeaders(page int, size int) ([]*types.Header, error) {

	offset := (page - 1) * size

	lastHeader, err := bc.ReadLastHeader()
	if err != nil {
		return nil, err
	}

	headers := []*types.Header{}

	if lastHeader == nil {
		return headers, nil
	}

	lastHeaderHeight := int(lastHeader.Height)
	start := lastHeaderHeight - offset
	if start < 0 {
		return headers, nil
	}

	end := start - size
	if end < -1 {
		end = -1
	}

	for i := start; i > end; i-- {
		header, err := bc.ReadHeaderByHeight(int32(i))
		if err != nil {
			return nil, err
		}
		if header == nil {
			return nil, fmt.Errorf("header %d is nil", i)
		}
		headers = append(headers, header)
	}

	return headers, nil
}

func (bc *Blockchain) ReadLastHeader() (*types.Header, error) {
	header, err := bc.db.SelectLastHeader()
	if err != nil {
		return nil, err
	}

	return header, nil
}

func (bc *Blockchain) ReadLastHeaderHeight() (*int32, error) {
	header, err := bc.ReadLastHeader()
	if err != nil {
		return nil, err
	}

	if header == nil {
		height := int32(-1)
		return &height, nil
	}

	return &header.Height, nil
}

func (bc *Blockchain) ReadTxByHash(hash common.Hash) (*types.Transaction, error) {
	tx, err := bc.db.SelectTxByHash(hash)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (bc *Blockchain) ReadTxByNumber(number uint32) (*types.Transaction, error) {
	tx, err := bc.db.SelectTxByNumber(number)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (bc *Blockchain) ReadTxs(page int, size int) ([]*types.Transaction, error) {
	offset := (page - 1) * size

	lastTxNumber, err := bc.ReadLastTxNumber()
	if err != nil {
		return nil, err
	}

	txs := []*types.Transaction{}

	if lastTxNumber == nil {
		return txs, nil
	}

	start := int(*lastTxNumber) - offset
	if start < 0 {
		return txs, nil
	}

	end := start - size
	if end < -1 {
		end = -1
	}

	for i := start; i > end; i-- {
		tx, err := bc.ReadTxByNumber(uint32(i))
		if err != nil {
			return nil, err
		}
		if tx == nil {
			return nil, fmt.Errorf("tx is nil")
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

func (bc *Blockchain) ReadLastTx() (*types.Transaction, error) {
	tx, err := bc.db.SelectLastTx()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (bc *Blockchain) ReadLastTxNumber() (*uint32, error) {
	number, err := bc.db.SelectLastTxNumber()
	if err != nil {
		return nil, err
	}

	return number, nil
}

func (bc *Blockchain) ReadAccountByAddress(address common.Address) (*types.Account, error) {
	account, err := bc.db.SelectAccountByAddress(address)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (bc *Blockchain) ReadBalance(address common.Address) (*uint64, error) {
	account, err := bc.ReadAccountByAddress(address)
	if err != nil {
		return nil, err
	}

	if account == nil {
		bal := uint64(0)
		return &bal, nil
	}

	balance := account.Balance
	return &balance, nil
}
