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

func (bc *Blockchain) WriteBlockWithHeight(height int32, block *types.Block) error {
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

func (bc *Blockchain) WriteHeaderWithHash(hash common.Hash, header *types.Header) error {
	if err := bc.db.InsertHeaderWithHash(hash, header); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteHeaderWithHeight(height int32, header *types.Header) error {
	if err := bc.db.InsertHeaderWithHeight(height, header); err != nil {
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

func (bc *Blockchain) WriteAccountWithAddress(address common.Address, account *types.Account) (*types.Account, error) {
	if account == nil {
		account = types.CreateAccount(address)
	}

	if err := bc.db.InsertAccountWithAddress(address, account); err != nil {
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

	if err = bc.db.InsertAccountWithAddress(fromAccount.Address, fromAccount); err != nil {
		return err
	}
	if err = bc.db.InsertAccountWithAddress(fromAccount.Address, fromAccount); err != nil {
		return err
	}
	return nil
}
