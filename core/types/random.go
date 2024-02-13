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
func NewRandomTransaction(privateKey crypto.PrivateKey) *Transaction {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return &Transaction{
		Nonce: 171, //ab
		From:  privateKey.PublicKey().Address(),
		To:    privateKey.PublicKey().Address(),
		Value: 171, //ab
		Data:  RandomBytes(r.Intn(1000)),
	}
}

func NewRandomTransactionWithSignature(t *testing.T, privateKey crypto.PrivateKey, size int) *Transaction {
	tx := NewRandomTransaction(privateKey)
	assert.Nil(t, tx.Sign(privateKey))
	return tx
}

func NewRandomBlock(t *testing.T, height int32, prevBlockHash common.Hash) *Block {
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

func NewRandomBlockWithSignature(t *testing.T, pk crypto.PrivateKey, height int32, prevHash common.Hash) *Block {
	b := NewRandomBlock(t, height, prevHash)
	assert.Nil(t, b.Sign(pk))

	return b
}
