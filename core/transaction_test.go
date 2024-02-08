package core

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/core/types"
	"testing"

	"github.com/barreleye-labs/barreleye/crypto"
	"github.com/stretchr/testify/assert"
)

func TestNativeTransferTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx := &types.Transaction{
		To:    toPrivKey.PublicKey().Address(),
		Value: 666,
	}

	assert.Nil(t, tx.Sign(fromPrivKey))
}

func TestSignTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &types.Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	tx := &types.Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.Nil(t, tx.Verify())

	otherPrivKey := crypto.GeneratePrivateKey()
	tx.Signer = otherPrivKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(types.NewGobTxEncoder(buf)))

	txDecoded := new(types.Transaction)
	assert.Nil(t, txDecoded.Decode(types.NewGobTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func randomTxWithSignature(t *testing.T) *types.Transaction {
	privateKey := crypto.GeneratePrivateKey()

	toPrivateKey := crypto.GeneratePrivateKey()
	toPublicKey := toPrivateKey.PublicKey()

	tx := types.Transaction{
		Data: []byte("foo"),
		To:   toPublicKey.Address(),
	}
	assert.Nil(t, tx.Sign(privateKey))

	return &tx
}
