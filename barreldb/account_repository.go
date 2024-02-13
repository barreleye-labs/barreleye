package barreldb

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (barrelDB *BarrelDatabase) InsertAccountWithAddress(address common.Address, account *types.Account) error {
	buf := &bytes.Buffer{}
	if err := account.Encode(types.NewGobAccountEncoder(buf)); err != nil {
		return err
	}

	if err := barrelDB.GetTable(AddressAccountTableName).Put(address.ToSlice(), buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func (barrelDB *BarrelDatabase) SelectAccountByAddress(address common.Address) (*types.Account, error) {
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
