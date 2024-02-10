package barreldb

import (
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
)

func (barrelDB *BarrelDatabase) InsertAccountWithAddress(address common.Address, account *types.Account) error {
	//buf := &bytes.Buffer{}
	//if err := block.Encode(types.NewGobBlockEncoder(buf)); err != nil {
	//	return err
	//}
	//
	//if err := barrelDB.GetTable(HashBlockTableName).Put(hash.ToSlice(), buf.Bytes()); err != nil {
	//	return err
	//}
	return nil
}
