package core

import (
	"github.com/barreleye-labs/barreleye/barreldb"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core/types"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddBlockToHeight(t *testing.T) {
	_ = barreldb.RemoveData("data")

	bc := newBlockchainWithGenesis(t)
	defer bc.db.Close()

	assert.Nil(t, bc.LinkBlock(randomBlock(t, 1, getPrevBlockHash(t, bc, int32(1)))))
	assert.NotNil(t, bc.LinkBlock(randomBlock(t, 3, common.Hash{})))
}

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	pk := types.GeneratePrivateKey()
	bc, err := NewBlockchain(log.NewNopLogger(), pk)
	assert.Nil(t, err)

	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height int32) common.Hash {
	prevHeader, err := bc.ReadHeaderByHeight(height - 1)
	assert.Nil(t, err)
	return types.BlockHasher{}.Hash(prevHeader)
}

func TestAddBlock(t *testing.T) {
	_ = barreldb.RemoveData("data")

	bc := newBlockchainWithGenesis(t)
	defer bc.db.Close()

	lenBlocks := 1000
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, int32(i+1), getPrevBlockHash(t, bc, int32(i+1)))
		assert.Nil(t, bc.LinkBlock(block))
	}

	lastBlockHeight, _ := bc.ReadLastBlockHeight()

	assert.Equal(t, *lastBlockHeight, int32(lenBlocks))
	assert.Equal(t, *lastBlockHeight+1, lenBlocks)
	assert.NotNil(t, bc.LinkBlock(randomBlock(t, 10, common.Hash{})))
}

func TestNewBlockchain(t *testing.T) {
	_ = barreldb.RemoveData("data")

	bc := newBlockchainWithGenesis(t)
	defer bc.db.Close()

	lastBlockHeight, _ := bc.ReadLastBlockHeight()

	assert.NotNil(t, bc.validator)
	assert.Equal(t, lastBlockHeight, int32(0))
}

func TestHasBlock(t *testing.T) {
	_ = barreldb.RemoveData("data")

	bc := newBlockchainWithGenesis(t)
	defer bc.db.Close()

	block, err := bc.ReadBlockByHeight(0)
	assert.Nil(t, err)

	assert.NotNil(t, block)

	block, err = bc.ReadBlockByHeight(1)
	assert.Nil(t, err)

	assert.Nil(t, block)

	block, err = bc.ReadBlockByHeight(100)
	assert.Nil(t, err)

	assert.Nil(t, block)
}
