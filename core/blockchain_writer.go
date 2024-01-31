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

func (bc *Blockchain) WriteTxWithHash(hash common.Hash, tx *types.Transaction) error {
	if err := bc.db.InsertTxWithHash(hash, tx); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteTxWithNumber(number uint32, tx *types.Transaction) error {
	if err := bc.db.InsertTxWithNumber(number, tx); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteLastTx(tx *types.Transaction) error {
	if err := bc.db.InsertLastTx(tx); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteLastTxNumber(number uint32) error {
	if err := bc.db.InsertLastTxNumber(number); err != nil {
		return err
	}
	return nil
}
