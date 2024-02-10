package types

import "github.com/barreleye-labs/barreleye/common"

type Account struct {
	Address common.Address
	Balance uint64
}

func CreateAccount(address common.Address) *Account {

	return nil
}
