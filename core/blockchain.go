package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/barreldb"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"sync"
	"time"

	"github.com/go-kit/log"
)

type Blockchain struct {
	logger    log.Logger
	lock      sync.RWMutex
	validator Validator
	db        *barreldb.BarrelDatabase
}

func NewBlockchain(l log.Logger, privateKey *types.PrivateKey) (*Blockchain, error) {
	db, _ := barreldb.New()

	if err := setTables(db); err != nil {
		return nil, err
	}

	bc := &Blockchain{
		logger: l,
		db:     db,
	}
	bc.validator = NewBlockValidator(bc)

	publicKey := privateKey.PublicKey
	coinbase, err := bc.ReadAccountByAddress(publicKey.Address())
	if err != nil {
		return nil, err
	}

	if coinbase == nil {
		coinbase = types.CreateAccount(publicKey.Address())
		if err := bc.WriteAccountWithAddress(publicKey.Address(), coinbase); err != nil {
			return nil, err
		}
	}

	if common.GetFlag("role") == "genesis" {
		lastBlock, err := bc.ReadLastBlock()
		if err != nil {
			return nil, err
		}

		if lastBlock == nil {
			err = bc.LinkBlockWithoutValidation(CreateGenesisBlock(privateKey))
			if err != nil {
				return nil, err
			}

			_ = bc.logger.Log("msg", "🌞 create genesis block")
		}
	}

	return bc, nil
}

func setTables(db *barreldb.BarrelDatabase) error {
	err := db.CreateTable(barreldb.HashBlockTableName, barreldb.HashBlockPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.HeightBlockTableName, barreldb.HeightBlockPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.LastBlockTableName, barreldb.LastBlockPrefix)
	if err != nil {
		return err
	}

	err = db.CreateTable(barreldb.HashHeaderTableName, barreldb.HashHeaderPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.HeightHeaderTableName, barreldb.HeightHeaderPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.LastHeaderTableName, barreldb.LastHeaderPrefix)
	if err != nil {
		return err
	}

	err = db.CreateTable(barreldb.HashTxTableName, barreldb.HashTxPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.NumberTxTableName, barreldb.NumberTxPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.LastTxTableName, barreldb.LastTxPrefix)
	if err != nil {
		return err
	}
	err = db.CreateTable(barreldb.LastTxNumberTableName, barreldb.LastTxNumberPrefix)
	if err != nil {
		return err
	}

	err = db.CreateTable(barreldb.AddressAccountTableName, barreldb.AddressAccountPrefix)
	if err != nil {
		return err
	}
	return nil
}

func CreateGenesisBlock(privateKey *types.PrivateKey) *types.Block {
	header := &types.Header{
		Version:   1,
		Height:    0,
		Timestamp: time.Now().UnixNano(),
	}

	b, _ := types.NewBlock(header, nil)

	if err := b.Sign(*privateKey); err != nil {
		panic(err)
	}
	return b
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) LinkBlock(b *types.Block) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	return bc.LinkBlockWithoutValidation(b)
}

func (bc *Blockchain) handleTransaction(tx *types.Transaction) error {
	if tx.From.Equal(tx.To) {
		return fmt.Errorf("from and to must be different")
	}

	fromAccount, err := bc.ReadAccountByAddress(tx.From)
	if err != nil {
		return err
	}

	if fromAccount == nil {
		fromAccount = types.CreateAccount(tx.From)
		if err = bc.WriteAccountWithAddress(fromAccount.Address, fromAccount); err != nil {
			return err
		}
	}

	if fromAccount.Nonce != tx.Nonce {
		return fmt.Errorf("invalid tx nonce")
	}

	toAccount, err := bc.ReadAccountByAddress(tx.To)
	if err != nil {
		return err
	}

	if toAccount == nil {
		toAccount = types.CreateAccount(tx.To)
		if err = bc.WriteAccountWithAddress(toAccount.Address, toAccount); err != nil {
			return err
		}
	}

	if err = fromAccount.Transfer(toAccount, tx.Value); err != nil {
		return err
	}

	fromAccount.Nonce++

	if err = bc.WriteAccountWithAddress(fromAccount.Address, fromAccount); err != nil {
		return err
	}
	if err = bc.WriteAccountWithAddress(toAccount.Address, toAccount); err != nil {
		return err
	}

	_ = bc.logger.Log(
		"msg", "transfer",
		"from", fromAccount.Address,
		"to", toAccount.Address,
		"value", tx.Value)

	return nil
}

func (bc *Blockchain) LinkBlockWithoutValidation(b *types.Block) error {
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

	/*	check sync account status
		barreleyeKey, _ := types.CreatePrivateKey("a2288db63c7016b815c55c1084c2491b8599834500408ba863ec379895373ae9")
		barreleye, _ := bc.ReadAccountByAddress(barreleyeKey.PublicKey.Address())
		fmt.Println("barreleye: ", barreleye)
		nayoungKey, _ := types.CreatePrivateKey("c4e0f3f39c5438d2f7ba8b830f5a5538c6a63c752cb36fb1b91911539af01421")
		nayoung, _ := bc.ReadAccountByAddress(nayoungKey.PublicKey.Address())
		fmt.Println("nayoung: ", nayoung)
		youngminKey, _ := types.CreatePrivateKey("f2e1e4331b10c2b84a8ed58226398f5d11ee78052afa641d16851bd66bbdadb7")
		youngmin, _ := bc.ReadAccountByAddress(youngminKey.PublicKey.Address())
		fmt.Println("youngmin: ", youngmin)
	*/

	_ = bc.logger.Log(
		"msg", "🔗 link new block",
		"hash", b.GetHash(),
		"height", b.Height,
		"txCount", len(b.Transactions),
	)

	return nil
}
