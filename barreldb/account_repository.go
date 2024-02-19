package barreldb

import (
	"bytes"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

// AddressAccount Repository
func (barrelDB *BarrelDatabase) UpsertAddressAccount(address common.Address, account *types.Account) error {
	buf := &bytes.Buffer{}
	if err := account.Encode(types.NewGobAccountEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(AddressAccountTableName).Put(address.ToSlice(), buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (barrelDB *BarrelDatabase) DeleteAddressAccount(address common.Address) error {
	if err := barrelDB.GetTable(AddressAccountTableName).Delete(address.ToSlice()); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) SelectAddressAccount(address common.Address) (*types.Account, error) {
	data, err := barrelDB.GetTable(AddressAccountTableName).Get(address.ToSlice())
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	account := new(types.Account)
	err = account.Decode(types.NewGobAccountDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (barrelDB *BarrelDatabase) SelectAccountBalance(address common.Address) (*uint64, error) {
	data, err := barrelDB.GetTable(AddressAccountTableName).Get(address.ToSlice())
	if err != nil {
		if err.Error() != common.LevelDBNotFoundError {
			return nil, err
		}
		return nil, nil
	}

	account := new(types.Account)
	err = account.Decode(types.NewGobAccountDecoder(bytes.NewBuffer(data)))
	if err != nil {
		return nil, err
	}

	return &account.Balance, nil
}

func (barrelDB *BarrelDatabase) IncreaseAccountBalance(address common.Address, amount uint64) error {
	account, err := barrelDB.SelectAddressAccount(address)
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("not found account")
	}

	account.Balance += amount
	if err = barrelDB.UpsertAddressAccount(account.Address, account); err != nil {
		return err
	}
	return nil
}

func (barrelDB *BarrelDatabase) DecreaseAccountBalance(address common.Address, amount uint64) error {
	account, err := barrelDB.SelectAddressAccount(address)
	if err != nil {
		return err
	}

	if account == nil {
		return fmt.Errorf("not found account")
	}

	if account.Balance < amount {
		return fmt.Errorf("not enough balance")
	}

	account.Balance -= amount
	if err = barrelDB.UpsertAddressAccount(account.Address, account); err != nil {
		return err
	}
	return nil
}
