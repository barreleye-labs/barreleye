package core

import (
	"fmt"
	"sync"

	"github.com/barreleye-labs/barreleye/types"
	"github.com/go-kit/log"
)

type Blockchain struct {
	logger     log.Logger
	store      Storage
	lock       sync.RWMutex
	headers    []*Header
	blocks     []*Block
	txStore    map[types.Hash]*Transaction
	blockStore map[types.Hash]*Block

	AccountState *AccountState

	stateLock       sync.RWMutex
	collectionState map[types.Hash]*CollectionTx
	minState        map[types.Hash]*MintTx
	validator       Validator
	// TODO: make this on interface.
	contractState *State
}

func NewBlockchain(l log.Logger, genesis *Block, ac *AccountState) (*Blockchain, error) {
	bc := &Blockchain{
		contractState:   NewState(),
		headers:         []*Header{},
		store:           NewMemorystore(),
		logger:          l,
		AccountState:    ac,
		collectionState: make(map[types.Hash]*CollectionTx),
		minState:        make(map[types.Hash]*MintTx),
		blockStore:      make(map[types.Hash]*Block),
		txStore:         make(map[types.Hash]*Transaction),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)

	return bc, err
}

func (bc *Blockchain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *Blockchain) AddBlock(b *Block) error {
	if err := bc.validator.ValidateBlock(b); err != nil {
		return err
	}

	bc.stateLock.Lock()
	defer bc.stateLock.Unlock()

	for _, tx := range b.Transactions {
		if len(tx.Data) > 0 {
			bc.logger.Log("msg", "excuting code", "len", len(tx.Data), "hash", tx.Hash(&TxHasher{}))
			vm := NewVM(tx.Data, bc.contractState)
			if err := vm.Run(); err != nil {
				return err
			}
		}

		// If the txInner of the transaction is not nil we need to handle
		//	the native NFT implementation.
		if tx.TxInner != nil {
			if err := bc.handleNativeNFT(tx); err != nil {
				return err
			}
		}

		// handle the native transaction here
		if tx.Value > 0 {
			if err := bc.handleNativeTransfter(tx); err != nil {
				return nil
			}
		}
	}

	fmt.Printf("%+v\n", bc.AccountState)

	return bc.addBlockWithoutValidation(b)
}

func (bc *Blockchain) handleNativeTransfter(tx *Transaction) error {
	bc.logger.Log("msg", "handle native token transfer", "from", tx.From, "to", tx.To, "value", tx.Value)

	return bc.AccountState.Transfer(tx.From.Address(), tx.To.Address(), tx.Value)
}

func (bc *Blockchain) handleNativeNFT(tx *Transaction) error {
	hash := tx.Hash(TxHasher{})
	switch t := tx.TxInner.(type) {
	case CollectionTx:
		bc.collectionState[hash] = &t
		bc.logger.Log("msg", "created new NFT collection", "hash", hash)
	case MintTx:
		_, ok := bc.collectionState[t.Collection]
		if !ok {
			return fmt.Errorf("collection (%s) does not exist on the blockchain", t.Collection)
		}

		bc.minState[hash] = &t

		bc.logger.Log("msg", "created new NFT mint", "NFT", t.NFT, "collection", t.Collection)
	default:
		fmt.Printf("unsupported tx type %v", t)
	}

	return nil
}

func (bc *Blockchain) GetBlockByHash(hash types.Hash) (*Block, error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	block, ok := bc.blockStore[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash (%s) not found", hash)
	}

	return block, nil
}

func (bc *Blockchain) GetBlock(height uint32) (*Block, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.blocks[height], nil
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("given height (%d) too high", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) GetTxByHash(hash types.Hash) (*Transaction, error) {
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

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	bc.blocks = append(bc.blocks, b)
	bc.blockStore[b.Hash(BlockHasher{})] = b

	for _, tx := range b.Transactions {
		bc.txStore[tx.Hash(TxHasher{})] = tx
	}

	bc.lock.Unlock()

	bc.logger.Log(
		"msg", "new block",
		"hash", b.Hash(BlockHasher{}),
		"height", b.Height,
		"transactions", len(b.Transactions),
	)

	return bc.store.Put(b)
}
