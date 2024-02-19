package types

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
)

type Account struct {
	Address common.Address
	Nonce   uint64
	Balance uint64
}

func CreateAccount(address common.Address) *Account {
	return &Account{
		Address: address,
		Nonce:   uint64(0),
		Balance: uint64(0),
	}
}

func (a *Account) Decode(dec Decoder[*Account]) error {
	return dec.Decode(a)
}

func (a *Account) Encode(enc Encoder[*Account]) error {
	return enc.Encode(a)
}

func (a *Account) Transfer(to *Account, amount uint64) error {
	if a.Balance < amount {
		return fmt.Errorf("insufficient account balance")
	}

	if err := a.SubBalance(amount); err != nil {
		return err
	}
	to.AddBalance(amount)
	return nil
}

func (a *Account) AddBalance(amount uint64) {
	a.Balance += amount
}

func (a *Account) SubBalance(amount uint64) error {
	if a.Balance < amount {
		return fmt.Errorf("balance cannot be negative")
	}
	a.Balance -= amount
	return nil
}
