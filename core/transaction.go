package core

import (
	"encoding/gob"
	"fmt"
	"math/rand"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/barreleye-labs/barreleye/types"
)

type TxType byte

const (
	TxTypeCollection TxType = iota // 0x0
	TxTypeMint					   // 0x01 
)

type CollectionTx struct {
	Fee		 int64
	MetaData []byte
}

type MintTx struct {
	Fee 	   		int64
	NFT 	   		types.Hash
	Collection 		types.Hash
	MetaData   		[]byte
	CollectionOwner crypto.PublicKey
	Signature		crypto.Signature
}

type Transaction struct {
	Type	  TxType
	TxInner	  any
	Data 	  []byte
	From      crypto.PublicKey
	Signature *crypto.Signature
	Nounce	  int64

	// cached version of the tx data hash
	hash types.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
		Nounce: rand.Int63n(1000000000000000),
	}
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	hash := tx.Hash(TxHasher{})
	sig, err := privKey.Sign(hash.ToSlice())
	if err != nil {
		return err
	}

	tx.From = privKey.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	hash := tx.Hash(TxHasher{})
	if !tx.Signature.Verify(tx.From, hash.ToSlice()) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(enc Encoder[*Transaction]) error {
	return enc.Encode(tx)
}

func init() {
	gob.Register(CollectionTx{})
	gob.Register(MintTx{})
}