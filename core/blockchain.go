package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/barreldb"
	"github.com/barreleye-labs/barreleye/core/types"
	"sync"
	"time"

	"github.com/go-kit/log"
)

type Blockchain struct {
	logger        log.Logger
	store         Storage
	lock          sync.RWMutex
	stateLock     sync.RWMutex
	validator     Validator
	contractState *State
	db            *barreldb.BarrelDatabase
}

func NewBlockchain(l log.Logger, privateKey *types.PrivateKey, nodeID string) (*Blockchain, error) {
	db, _ := barreldb.New()

	err := db.CreateTable(barreldb.HashBlockTableName, barreldb.HashBlockPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.HeightBlockTableName, barreldb.HeightBlockPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastBlockTableName, barreldb.LastBlockPrefix)
	if err != nil {
		return nil, err
	}

	err = db.CreateTable(barreldb.HashHeaderTableName, barreldb.HashHeaderPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.HeightHeaderTableName, barreldb.HeightHeaderPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastHeaderTableName, barreldb.LastHeaderPrefix)
	if err != nil {
		return nil, err
	}

	err = db.CreateTable(barreldb.HashTxTableName, barreldb.HashTxPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.NumberTxTableName, barreldb.NumberTxPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastTxTableName, barreldb.LastTxPrefix)
	if err != nil {
		return nil, err
	}
	err = db.CreateTable(barreldb.LastTxNumberTableName, barreldb.LastTxNumberPrefix)
	if err != nil {
		return nil, err
	}

	err = db.CreateTable(barreldb.AddressAccountTableName, barreldb.AddressAccountPrefix)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		contractState: NewState(),
		store:         NewMemorystore(),
		logger:        l,
		db:            db,
	}
	bc.validator = NewBlockValidator(bc)

	if privateKey != nil {
		coinbase := privateKey.PublicKey

		coinbaseAccount := types.CreateAccount(coinbase.Address())
		if err = bc.WriteAccountWithAddress(coinbase.Address(), coinbaseAccount); err != nil {
			return nil, err
		}
	}

	if nodeID == "GENESIS-NODE" {
		lastBlock, err := bc.ReadLastBlock()
		if err != nil {
			return nil, err
		}

		if lastBlock == nil {
			err = bc.addBlockWithoutValidation(CreateGenesisBlock(privateKey))
			if err != nil {
				return nil, err
			}

			_ = bc.logger.Log("msg", "🌞 create genesis block")
		}
	}

	return bc, nil
}

func CreateGenesisBlock(privateKey *types.PrivateKey) *types.Block {
	//coinbase := privateKey.PublicKey()

	//tx := &types.Transaction{
	//	Nonce: 171, //ab
	//	From:  coinbase.Address(),
	//	To:    coinbase.Address(),
	//	Value: 171, //ab
	//	Data:  []byte{171},
	//}

	//if err := tx.Sign(*privateKey); err != nil {
	//	panic(err)
	//}

	header := &types.Header{
		Version:   1,
		Height:    0,
		Timestamp: time.Now().UnixNano(),
	}

	b, _ := types.NewBlock(header, nil)

	//b.Transactions = append(b.Transactions, tx)
	//hash, _ := types.CalculateDataHash(b.Transactions)
	//b.DataHash = hash

	if err := b.Sign(*privateKey); err != nil {
		panic(err)
	}
	return b
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *types.Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) handleTransaction(tx *types.Transaction) error {
	account, err := bc.ReadAccountByAddress(tx.From)
	if err != nil {
		return err
	}

	if account == nil {
		account = types.CreateAccount(tx.From)
		if err = bc.WriteAccountWithAddress(account.Address, account); err != nil {
			return err
		}
	}

	if account.Nonce != tx.Nonce {
		return fmt.Errorf("invalid tx nonce")
	}

	if len(tx.Data) > 0 {
		_ = bc.logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.GetHash())

		vm := NewVM(tx.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
	}

	if err := bc.Transfer(tx.From, tx.To, tx.Value); err != nil {
		return err
	}

	account.Nonce++

	if err = bc.WriteAccountWithAddress(account.Address, account); err != nil {
		return err
	}

	return nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *types.Block) error {
	for i := 0; i < len(b.Transactions); i++ {
		if err := bc.handleTransaction(b.Transactions[i]); err != nil {
			_ = bc.logger.Log("error", err.Error())

			b.Transactions[i] = b.Transactions[len(b.Transactions)-1]
			b.Transactions = b.Transactions[:len(b.Transactions)-1]
			i--
		}
	}

	if err := bc.WriteBlockWithHash(b.GetHash(), b); err != nil {
		return err
	}
	if err := bc.WriteBlockWithHeight(b.Height, b); err != nil {
		return err
	}
	if err := bc.WriteLastBlock(b); err != nil {
		return err
	}

	if err := bc.WriteHeaderWithHash(b.GetHash(), b.Header); err != nil {
		return err
	}
	if err := bc.WriteHeaderWithHeight(b.Height, b.Header); err != nil {
		return err
	}
	if err := bc.WriteLastHeader(b.Header); err != nil {
		return err
	}

	if err := bc.GiveReward(b.Signer.Address()); err != nil {
		return err
	}

	for _, tx := range b.Transactions {
		nextTxNumber := uint32(0)
		lastTxNumber, err := bc.ReadLastTxNumber()
		if err != nil {
			return err
		}

		if lastTxNumber != nil {
			nextTxNumber = *lastTxNumber + 1
		}

		if err = bc.WriteTxWithHash(tx.GetHash(), tx); err != nil {
			return err
		}
		if err = bc.WriteTxWithNumber(nextTxNumber, tx); err != nil {
			return err
		}
		if err = bc.WriteLastTx(tx); err != nil {
			return err
		}
		if err = bc.WriteLastTxNumber(nextTxNumber); err != nil {
			return err
		}
	}

	_ = bc.logger.Log(
		"msg", "🔗 add new block",
		"hash", b.GetHash(),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return nil
}
