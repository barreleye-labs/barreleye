package types

import (
	"github.com/barreleye-labs/barreleye/common"
	"math/rand"
	"testing"
	"time"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/stretchr/testify/assert"
)

func RandomBytes(size int) []byte {
	token := make([]byte, size)
	rand.Read(token)
	return token
}

func RandomHash() common.Hash {
	return common.HashFromBytes(RandomBytes(32))
}

// NewRandomTransaction return a new random transaction whithout signature.
func NewRandomTransaction(size int) *Transaction {
	return NewTransaction(RandomBytes(size))
}

func NewRandomTransactionWithSignature(t *testing.T, privKey crypto.PrivateKey, size int) *Transaction {
	tx := NewRandomTransaction(size)
	assert.Nil(t, tx.Sign(privKey))
	return tx
}

func NewRandomBlock(t *testing.T, height uint32, prevBlockHash common.Hash) *Block {
	txSigner := crypto.GeneratePrivateKey()
	tx := NewRandomTransactionWithSignature(t, txSigner, 100)
	header := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}
	b, err := NewBlock(header, []*Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash

	return b
}

func NewRandomBlockWithSignature(t *testing.T, pk crypto.PrivateKey, height uint32, prevHash common.Hash) *Block {
	b := NewRandomBlock(t, height, prevHash)
	assert.Nil(t, b.Sign(pk))

	return b
}
