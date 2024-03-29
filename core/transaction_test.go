package core

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/core/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransfer(t *testing.T) {
	fromPrivateKey := types.GeneratePrivateKey()
	toPrivateKey := types.GeneratePrivateKey()
	tx := &types.Transaction{
		To:    toPrivateKey.PublicKey.Address(),
		Value: 666,
	}

	assert.Nil(t, tx.Sign(fromPrivateKey))
}

func TestSignTx(t *testing.T) {
	privKey := types.GeneratePrivateKey()
	tx := &types.Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTx(t *testing.T) {
	signerPrivateKey := types.GeneratePrivateKey()
	tx := &types.Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(signerPrivateKey))
	assert.Nil(t, tx.Verify())

	hackerPrivateKey := types.GeneratePrivateKey()
	tx.Signer = hackerPrivateKey.PublicKey

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
	privateKey := types.GeneratePrivateKey()

	toPrivateKey := types.GeneratePrivateKey()
	toPublicKey := toPrivateKey.PublicKey

	tx := types.Transaction{
		Nonce:  171, //ab
		From:   privateKey.PublicKey.Address(),
		To:     toPublicKey.Address(),
		Value:  171,         //ab
		Data:   []byte{171}, //ab
		Signer: privateKey.PublicKey,
	}
	assert.Nil(t, tx.Sign(privateKey))

	return &tx
}
