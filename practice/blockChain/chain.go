package blockchain

import (
	"fmt"

	"github.com/dgraph-io/badger"
)

const (
	dbpath = "./tmp/blocks"
)

var key = []byte("lh")

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbpath)
	opts.Dir = dbpath
	opts.ValueDir = dbpath

	db, err := badger.Open(opts)

	ErrorHandle(err)

	err = db.Update(func(tx *badger.Txn) error {
		if _, err := tx.Get(key); err == badger.ErrKeyNotFound {
			fmt.Println("no exist")
			genesis := Genesis()
			err = tx.Set(genesis.Hash, genesis.Serialize())
			err = tx.Set(key, genesis.Hash)

			lastHash = genesis.Hash

			return err
		} else {
			item, err := tx.Get(key)
			ErrorHandle(err)
			lastHash, err = item.ValueCopy(key)
			return err
		}
	})
	ErrorHandle(err)

	return &BlockChain{LastHash: lastHash, Database: db}
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte
	err := chain.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(key)

		ErrorHandle(err)

		lastHash, err = item.ValueCopy(key)

		return err
	})

	ErrorHandle(err)

	newBlock := CreateBlock(data, lastHash)

	err = chain.Database.Update(func(tx *badger.Txn) error {
		err := tx.Set(newBlock.Hash, newBlock.Serialize())

		ErrorHandle(err)

		err = tx.Set(key, newBlock.Hash)

		chain.LastHash = newBlock.Hash
		return err
	})

	ErrorHandle(err)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{chain.LastHash, chain.Database}
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(iter.CurrentHash)

		ErrorHandle(err)
		encodedBlock, err := item.ValueCopy(key)

		block = Deserialize(encodedBlock)

		return err
	})

	ErrorHandle(err)

	return block
}

func (iter *BlockChainIterator) GetByPrevHash(prevHash []byte) *Block {
	var block *Block

	err := iter.Database.View(func(tx *badger.Txn) error {
		item, err := tx.Get(prevHash)

		if err != nil {
			return nil
		}

		encodedBlock, err := item.ValueCopy(key)

		block = Deserialize(encodedBlock)

		return err
	})

	if err != nil {
		return nil
	}

	return block
}
