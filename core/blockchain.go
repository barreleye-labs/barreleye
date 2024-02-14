package core

import (
	"github.com/barreleye-labs/barreleye/barreldb"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/barreleye-labs/barreleye/crypto"
	"sync"

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

func NewBlockchain(l log.Logger, privateKey *crypto.PrivateKey, genesis *types.Block) (*Blockchain, error) {
	accountState := NewAccountState()

	if genesis != nil {
		coinbase := privateKey.PublicKey()
		accountState.CreateAccount(coinbase.Address())
	}

	db, _ := barreldb.New()

	/*
		bc 객체가 없는 영역에서 db 활용 Sample

		_ = db.CreateTable("block", barreldb.BlockPrefix)
		_ = db.CreateBlock("kim", "youngmin")
		data, _ := db.GetBlock("kim")
	*/

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
		coinbase := privateKey.PublicKey()
		accountState.CreateAccount(coinbase.Address())

		coinbaseAccount := types.CreateAccount(coinbase.Address())
		_, err = bc.WriteAccountWithAddress(coinbase.Address(), coinbaseAccount)
		if err != nil {
			return nil, err
		}
	}

	if genesis != nil {
		err = bc.addBlockWithoutValidation(genesis)
		if err != nil {
			return nil, err
		}
	}

	return bc, nil
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *types.Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) handleTransaction(tx *types.Transaction) error {
	if len(tx.Data) > 0 {
		bc.logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.GetHash())

		vm := NewVM(tx.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
	}

	if err := bc.Transfer(tx.From, tx.To, tx.Value); err != nil {
		return err
	}
	return nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *types.Block) error {
	bc.stateLock.Lock()
	for i := 0; i < len(b.Transactions); i++ {
		if err := bc.handleTransaction(b.Transactions[i]); err != nil {
			_ = bc.logger.Log("error", err.Error())

			b.Transactions[i] = b.Transactions[len(b.Transactions)-1]
			b.Transactions = b.Transactions[:len(b.Transactions)-1]

			continue
		}
	}
	bc.stateLock.Unlock()

	// fmt.Println("========ACCOUNT STATE==============")
	// fmt.Printf("%+v\n", bc.accountState.accounts)
	// fmt.Println("========ACCOUNT STATE==============")

	bc.lock.Lock()

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

	//data, _ := bc.ReadBlockByHash(b.GetHash(types.BlockHasher{}))
	//val := hexutil.Encode(data.Hash.ToSlice())
	//fmt.Println("hashblock::: ", val)
	//fmt.Println("fefefefefk::: ", data.Hash.String())

	//data, _ = bc.ReadBlockByHeight(b.Height)
	//fmt.Println("heightblock::: ", data)

	//data, _ = bc.ReadLastBlock()
	//fmt.Println("Lastblock::: ", data)

	for _, tx := range b.Transactions {
		nextTxNumber := uint32(0)
		number, err := bc.ReadLastTxNumber()
		if err != nil {
			return err
		}

		if number != nil {
			nextTxNumber = *number + 1
		}

		if err := bc.WriteTxWithHash(tx.GetHash(), tx); err != nil {
			return err
		}
		if err := bc.WriteTxWithNumber(nextTxNumber, tx); err != nil {
			return err
		}
		if err := bc.WriteLastTx(tx); err != nil {
			return err
		}
		if err := bc.WriteLastTxNumber(nextTxNumber); err != nil {
			return err
		}
	}
	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "🔗 add new block",
		"hash", b.GetHash(),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
