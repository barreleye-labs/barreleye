package types

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"math/rand"

	"github.com/barreleye-labs/barreleye/crypto"
)

type Transaction struct {
	Nonce     int64
	From      common.Address
	To        common.Address
	Value     uint64
	Data      []byte
	Signer    crypto.PublicKey
	Signature *crypto.Signature

	Hash common.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:  data,
		Nonce: rand.Int63n(1000000000000000),
	}
}

func (tx *Transaction) GetHash() common.Hash {
	if tx.Hash.IsZero() {
		tx.Hash = TxHasher{}.Hash(tx)
	}
	return tx.Hash
}

func (tx *Transaction) Sign(privateKey crypto.PrivateKey) error {
	hash := tx.GetHash()
	sig, err := privateKey.Sign(hash.ToSlice())
	if err != nil {
		return err
	}

	tx.Signer = privateKey.PublicKey()
	tx.Signature = sig

	return nil
}

func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	hash := tx.GetHash()
	if !tx.Signature.Verify(tx.Signer, hash.ToSlice()) {
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
