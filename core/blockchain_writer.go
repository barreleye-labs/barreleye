package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/config"
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
	if err := bc.db.UpsertLastTx(tx); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteLastTxNumber(number uint32) error {
	if err := bc.db.UpsertLastTxNumber(number); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) WriteAccountWithAddress(address common.Address, account *types.Account) error {
	if err := bc.db.UpsertAddressAccount(address, account); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) Transfer(from, to common.Address, amount uint64) error {
	fromAccount, err := bc.ReadAccountByAddress(from)
	if err != nil {
		return err
	}

	if fromAccount == nil {
		fromAccount = types.CreateAccount(from)
		if err = bc.WriteAccountWithAddress(from, fromAccount); err != nil {
			return err
		}
	}

	toAccount, err := bc.ReadAccountByAddress(to)
	if err != nil {
		return err
	}

	if toAccount == nil {
		toAccount = types.CreateAccount(to)
		if err = bc.WriteAccountWithAddress(to, toAccount); err != nil {
			return err
		}
	}

	if err = fromAccount.Transfer(toAccount, amount); err != nil {
		return err
	}

	if err = bc.db.UpsertAddressAccount(fromAccount.Address, fromAccount); err != nil {
		return err
	}
	if err = bc.db.UpsertAddressAccount(fromAccount.Address, fromAccount); err != nil {
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
	if err := bc.removeLastHeader(); err != nil {
		return err
	}

	if err := bc.removeLastBlockTxs(); err != nil {
		return err
	}

	lastBlock, err := bc.db.SelectLastBlock()
	if err != nil {
		return err
	}

	if lastBlock == nil {
		return fmt.Errorf("not found last block for removing block")
	}

	if lastBlock.Height < 1 {
		return fmt.Errorf("genesis block can not delete")
	}

	if err = bc.cancelReward(lastBlock.Validator.Address()); err != nil {
		return err
	}

	if err = bc.db.DeleteHashBlock(lastBlock.Hash); err != nil {
		return err
	}
	if err = bc.db.DeleteHeightBlock(lastBlock.Height); err != nil {
		return err
	}
	if err = bc.db.DeleteLastBlock(); err != nil {
		return err
	}

	prevBlock, err := bc.db.SelectHeightBlock(lastBlock.Height - 1)
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

func (bc *Blockchain) removeLastBlockTxs() error {
	lastBlock, err := bc.db.SelectLastBlock()
	if err != nil {
		return err
	}

	if lastBlock == nil {
		return fmt.Errorf("not found last block for removing txs")
	}

	lastTxNumber, err := bc.db.SelectLastTxNumber()
	if err != nil {
		return err
	}

	targetTxNum := uint32(0)
	if lastTxNumber != nil {
		targetTxNum = *lastTxNumber
	}

	edited := false
	for i := 0; i < len(lastBlock.Transactions); i++ {
		edited = true
		account, err := bc.db.SelectAddressAccount(lastBlock.Transactions[i].From)
		if err != nil {
			return err
		}

		if account == nil {
			return fmt.Errorf("not found tx from account")
		}

		account.Nonce--

		if err = bc.db.UpsertAddressAccount(account.Address, account); err != nil {
			return nil
		}

		if err = bc.db.DeleteHashTx(lastBlock.Transactions[i].Hash); err != nil {
			return err
		}

		if err = bc.db.DeleteNumberTx(targetTxNum); err != nil {
			return err
		}
		targetTxNum--
	}

	if edited {
		if err = bc.db.UpsertLastTxNumber(targetTxNum); err != nil {
			return err
		}

		lastTx, err := bc.db.SelectNumberTx(targetTxNum)
		if err != nil {
			return err
		}

		if lastTx == nil {
			return fmt.Errorf("not found numberTx")
		}

		if err = bc.db.UpsertLastTx(lastTx); err != nil {
			return err
		}
	}
	return nil
}

func (bc *Blockchain) removeLastHeader() error {
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

func (bc *Blockchain) cancelReward(address common.Address) error {
	account, err := bc.ReadAccountByAddress(address)
	if err != nil {
		return err
	}

	if account == nil {
		account = types.CreateAccount(address)
		if err = bc.WriteAccountWithAddress(account.Address, account); err != nil {
			return err
		}
	}

	if err = bc.db.DecreaseAccountBalance(address, config.BlockReward); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) GiveReward(address common.Address) error {
	account, err := bc.ReadAccountByAddress(address)
	if err != nil {
		return err
	}

	if account == nil {
		account = types.CreateAccount(address)
		if err = bc.WriteAccountWithAddress(account.Address, account); err != nil {
			return err
		}
	}

	if err = bc.db.IncreaseAccountBalance(account.Address, uint64(config.BlockReward)); err != nil {
		return err
	}
	return nil
}
