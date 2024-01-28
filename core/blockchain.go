package core

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/barreldb"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"sync"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/go-kit/log"
)

type Blockchain struct {
	logger log.Logger
	store  Storage
	// TODO: double check this!
	lock            sync.RWMutex
	headers         []*types.Header
	blocks          []*types.Block
	txStore         map[common.Hash]*types.Transaction
	blockStore      map[common.Hash]*types.Block
	accountState    *AccountState
	stateLock       sync.RWMutex
	collectionState map[common.Hash]*types.CollectionTx
	mintState       map[common.Hash]*types.MintTx
	validator       Validator
	// TODO: make this an interface.
	contractState *State
	db            *barreldb.BarrelDatabase
}

func NewBlockchain(l log.Logger, genesis *types.Block) (*Blockchain, error) {
	// We should create all states inside the scope of the newblockchain.
	// TODO: read this from disk later on
	accountState := NewAccountState()

	coinbase := crypto.PublicKey{}
	accountState.CreateAccount(coinbase.Address())

	db, _ := barreldb.New()

	/*
		bc 객체가 없는 영역에서 db 활용 Sample

		_ = db.CreateTable("block", barreldb.BlockPrefix)
		_ = db.CreateBlock("kim", "youngmin")
		data, _ := db.GetBlock("kim")
	*/

	err := db.CreateTable(barreldb.BlockTableName, barreldb.BlockPrefix)
	if err != nil {
		return nil, fmt.Errorf("fail to create table %s", barreldb.BlockTableName)
	}

	bc := &Blockchain{
		contractState:   NewState(),
		headers:         []*types.Header{},
		store:           NewMemorystore(),
		logger:          l,
		accountState:    accountState,
		collectionState: make(map[common.Hash]*types.CollectionTx),
		mintState:       make(map[common.Hash]*types.MintTx),
		blockStore:      make(map[common.Hash]*types.Block),
		txStore:         make(map[common.Hash]*types.Transaction),
		db:              db,
	}
	bc.validator = NewBlockValidator(bc)
	err = bc.addBlockWithoutValidation(genesis)

	//_ = bc.CreateBlock("kim", "youngmin")
	//data, _ := bc.GetBlockFromDB("kim")
	//fmt.Println("data::: ", data)

	return bc, err
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

func (bc *Blockchain) handleNativeTransfer(tx *types.Transaction) error {
	bc.logger.Log(
		"msg", "handle native token transfer",
		"from", tx.From,
		"to", tx.To,
		"value", tx.Value)

	return bc.accountState.Transfer(tx.From.Address(), tx.To.Address(), tx.Value)
}

func (bc *Blockchain) handleNativeNFT(tx *types.Transaction) error {
	hash := tx.GetHash(types.TxHasher{})

	switch t := tx.TxInner.(type) {
	case types.CollectionTx:
		bc.collectionState[hash] = &t
		bc.logger.Log("msg", "created new NFT collection", "hash", hash)
	case types.MintTx:
		_, ok := bc.collectionState[t.Collection]
		if !ok {
			return fmt.Errorf("collection (%s) does not exist on the blockchain", t.Collection)
		}
		bc.mintState[hash] = &t

		bc.logger.Log("msg", "created new NFT mint", "NFT", t.NFT, "collection", t.Collection)
	default:
		return fmt.Errorf("unsupported tx type %v", t)
	}

	return nil
}

func (bc *Blockchain) GetBlockByHash(hash common.Hash) (*types.Block, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	block, ok := bc.blockStore[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash (%s) not found", hash)
	}

	return block, nil
}

func (bc *Blockchain) GetBlock(height uint32) (*types.Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.blocks[height], nil
}

func (bc *Blockchain) GetHeader(height uint32) (*types.Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) GetTxByHash(hash common.Hash) (*types.Transaction, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	tx, ok := bc.txStore[hash]
	if !ok {
		return nil, fmt.Errorf("could not find tx with hash (%s)", hash)
	}

	return tx, nil
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

// [0, 1, 2 ,3] => 4 len
// [0, 1, 2 ,3] => 3 height
func (bc *Blockchain) Height() uint32 {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) handleTransaction(tx *types.Transaction) error {
	// If we have data inside execute that data on the VM.
	if len(tx.Data) > 0 {
		bc.logger.Log("msg", "executing code", "len", len(tx.Data), "hash", tx.GetHash(&types.TxHasher{}))

		vm := NewVM(tx.Data, bc.contractState)
		if err := vm.Run(); err != nil {
			return err
		}
	}

	// If the txInner of the transaction is not nil we need to handle
	// the native NFT implemtation.
	if tx.TxInner != nil {
		if err := bc.handleNativeNFT(tx); err != nil {
			return err
		}
	}

	// Handle the native transaction here
	if tx.Value > 0 {
		if err := bc.handleNativeTransfer(tx); err != nil {
			return err
		}
	}

	return nil
}

func (bc *Blockchain) addBlockWithoutValidation(b *types.Block) error {
	bc.stateLock.Lock()
	for i := 0; i < len(b.Transactions); i++ {
		if err := bc.handleTransaction(b.Transactions[i]); err != nil {
			bc.logger.Log("error", err.Error())

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
	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, b)

	_ = bc.CreateBlock(b.GetHash(types.BlockHasher{}), b)
	data, _ := bc.GetBlockByHashFromDB(b.GetHash(types.BlockHasher{}))
	fmt.Println("bbbb::: ", b)
	fmt.Println("data::: ", data)

	bc.blockStore[b.GetHash(types.BlockHasher{})] = b

	for _, tx := range b.Transactions {
		bc.txStore[tx.GetHash(types.TxHasher{})] = tx
	}
	bc.lock.Unlock()

	//bc.db.Put()
	bc.logger.Log(
		"msg", "🔗 add new block",
		"hash", b.GetHash(types.BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
