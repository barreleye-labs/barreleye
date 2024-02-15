package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (bc *Blockchain) WriteBlockWithHash(hash common.Hash, block *types.Block) error {
	if err := bc.db.InsertHashBlock(hash, block); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteBlockWithHeight(height int32, block *types.Block) error {
	if err := bc.db.InsertHeightBlock(height, block); err != nil {
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

func (bc *Blockchain) WriteHeaderWithHash(hash common.Hash, header *types.Header) error {
	if err := bc.db.InsertHashHeader(hash, header); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteHeaderWithHeight(height int32, header *types.Header) error {
	if err := bc.db.InsertHeightHeader(height, header); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteLastHeader(header *types.Header) error {
	if err := bc.db.InsertLastHeader(header); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteTxWithHash(hash common.Hash, tx *types.Transaction) error {
	if err := bc.db.InsertHashTx(hash, tx); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteTxWithNumber(number uint32, tx *types.Transaction) error {
	if err := bc.db.InsertNumberTx(number, tx); err != nil {
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

func (bc *Blockchain) WriteAccountWithAddress(address common.Address, account *types.Account) (*types.Account, error) {
	if account == nil {
		account = types.CreateAccount(address)
	}

	if err := bc.db.InsertAddressAccount(address, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (bc *Blockchain) Transfer(from, to common.Address, amount uint64) error {
	fromAccount, err := bc.ReadAccountByAddress(from)
	if err != nil {
		return err
	}

	if fromAccount == nil {
		// TODO: Register account at Coin Faucet
		fromAccount, err = bc.WriteAccountWithAddress(from, nil)
		if err != nil {
			return err
		}
	}

	toAccount, err := bc.ReadAccountByAddress(to)
	if err != nil {
		return err
	}

	if toAccount == nil {
		toAccount, err = bc.WriteAccountWithAddress(to, nil)
		if err != nil {
			return err
		}
	}

	if err = fromAccount.Transfer(toAccount, amount); err != nil {
		return err
	}

	if err = bc.db.InsertAddressAccount(fromAccount.Address, fromAccount); err != nil {
		return err
	}
	if err = bc.db.InsertAddressAccount(fromAccount.Address, fromAccount); err != nil {
		return err
	}

	_ = bc.logger.Log(
		"msg", "handle native token transfer",
		"from", fromAccount.Address,
		"to", toAccount.Address,
		"value", amount)
	return nil
}

func (bc *Blockchain) RemoveLastBlock() error {
	block, err := bc.db.SelectLastBlock()
	if err != nil {
		return err
	}

	if block == nil {
		return fmt.Errorf("not found last block for removing")
	}

	if block.Height < 1 {
		return fmt.Errorf("genesis block can not delete")
	}

	if err = bc.db.DeleteHashBlock(block.Hash); err != nil {
		return err
	}
	if err = bc.db.DeleteHeightBlock(block.Height); err != nil {
		return err
	}
	if err = bc.db.DeleteLastBlock(); err != nil {
		return err
	}

	prevBlock, err := bc.db.SelectHeightBlock(block.Height - 1)
	if err != nil {
		return err
	}

	if prevBlock == nil {
		return fmt.Errorf("not found previous block for inserting new last block")
	}

	if err = bc.db.InsertLastBlock(prevBlock); err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) RemoveLastHeader() error {
	header, err := bc.db.SelectLastHeader()
	if err != nil {
		return err
	}

	if header == nil {
		return fmt.Errorf("not found last header for removing")
	}

	if header.Height < 1 {
		return fmt.Errorf("genesis header can not delete")
	}

	if err = bc.db.DeleteHashHeader(types.BlockHasher{}.Hash(header)); err != nil {
		return err
	}
	if err = bc.db.DeleteHeightHeader(header.Height); err != nil {
		return err
	}
	if err = bc.db.DeleteLastHeader(); err != nil {
		return err
	}

	prevHeader, err := bc.db.SelectHeightHeader(header.Height - 1)
	if err != nil {
		return err
	}

	if prevHeader == nil {
		return fmt.Errorf("not found previous header for inserting new last header")
	}

	if err = bc.db.InsertLastHeader(prevHeader); err != nil {
		return err
	}

	return nil
}
