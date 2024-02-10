package core

import (
	"errors"
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"sync"
)

var (
	ErrAccountNotFound     = errors.New("account not found")
	ErrInsufficientBalance = errors.New("insufficient account balance")
)

type Account struct {
	Address common.Address
	Balance uint64
}

func (a *Account) String() string {
	return fmt.Sprintf("%d", a.Balance)
}

type AccountState struct {
	mu       sync.RWMutex
	accounts map[common.Address]*Account
}

func NewAccountState() *AccountState {
	return &AccountState{
		accounts: make(map[common.Address]*Account),
	}
}

func (s *AccountState) CreateAccount(address common.Address) *Account {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.CreateAccountWithoutLock(address)
}

func (s *AccountState) CreateAccountWithoutLock(address common.Address) *Account {
	acc := &Account{Address: address, Balance: 100_000_000_000}
	s.accounts[address] = acc
	return acc
}

func (s *AccountState) GetAccount(address common.Address) (*Account, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.getAccountWithoutLock(address)
}

func (s *AccountState) getAccountWithoutLock(address common.Address) (*Account, error) {
	account, ok := s.accounts[address]
	if !ok {
		return s.CreateAccountWithoutLock(address), nil
	}

	return account, nil
}

func (s *AccountState) GetBalance(address common.Address) (uint64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	account, err := s.getAccountWithoutLock(address)
	if err != nil {
		return 0, err
	}

	return account.Balance, nil
}

func (s *AccountState) Transfer(from, to common.Address, amount uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fromAccount, err := s.getAccountWithoutLock(from)
	if err != nil {
		return err
	}

	if fromAccount.Address.String() != "996fb92427ae41e4649b934ca495991b7852b855" {
		if fromAccount.Balance < amount {
			return ErrInsufficientBalance
		}
	}

	if fromAccount.Balance != 0 {
		fromAccount.Balance -= amount
	}

	if s.accounts[to] == nil {
		s.accounts[to] = &Account{
			Address: to,
		}
	}

	s.accounts[to].Balance += amount

	return nil
}
