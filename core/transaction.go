package core

import "github.com/barreleye-labs/barreleye/crypto"

type Transaction struct {
	Data []byte

	PublicKey crypto.PublicKey
	Signature *crypto.Signature
}

func (tx *Transaction) Sign(privKey crypto.PrivateKey) error {
	sig, err := privKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.PublicKey = privKey.PublicKey()
	tx.Signature = sig

	return nil
}
