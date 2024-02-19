package core

import (
	"bytes"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSignBlock(t *testing.T) {
	privKey := types.GeneratePrivateKey()
	b := randomBlock(t, 0, common.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privKey := types.GeneratePrivateKey()
	b := randomBlock(t, 0, common.Hash{})

	assert.Nil(t, b.Sign(privKey))
	assert.Nil(t, b.Verify())

	otherPrivKey := types.GeneratePrivateKey()
	b.Validator = otherPrivKey.PublicKey()
	assert.NotNil(t, b.Verify())

	b.Height = 100
	assert.NotNil(t, b.Verify())
}

func randomBlock(t *testing.T, height uint32, prevBlockHash common.Hash) *types.Block {
	privKey := types.GeneratePrivateKey()
	tx := randomTxWithSignature(t)
	header := &types.Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Height:        height,
		Timestamp:     time.Now().UnixNano(),
	}

	b, err := types.NewBlock(header, []*types.Transaction{tx})
	assert.Nil(t, err)
	dataHash, err := types.CalculateDataHash(b.Transactions)
	assert.Nil(t, err)
	b.Header.DataHash = dataHash
	assert.Nil(t, b.Sign(privKey))
	return b
}

func TestDecodeEncodeBlock(t *testing.T) {
	b := randomBlock(t, 1, common.Hash{})
	buf := &bytes.Buffer{}
	assert.Nil(t, b.Encode(types.NewGobBlockEncoder(buf)))

	bDecode := new(types.Block)
	assert.Nil(t, bDecode.Decode(types.NewGobBlockDecoder(buf)))

	assert.Equal(t, b.Header, bDecode.Header)

	for i := 0; i < len(b.Transactions); i++ {
		assert.Equal(t, b.Transactions[i], bDecode.Transactions[i])
	}

	assert.Equal(t, b.Validator, bDecode.Validator)
	assert.Equal(t, b.Signature, bDecode.Signature)
}
