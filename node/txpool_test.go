package node

import (
	"github.com/barreleye-labs/barreleye/core/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxSortedMapFirst(t *testing.T) {
	privateKey := types.GeneratePrivateKey()

	m := NewTxSortedMap()
	first := types.NewRandomTransaction(privateKey)
	m.Add(first)
	m.Add(types.NewRandomTransaction(privateKey))
	m.Add(types.NewRandomTransaction(privateKey))
	m.Add(types.NewRandomTransaction(privateKey))
	m.Add(types.NewRandomTransaction(privateKey))
	assert.Equal(t, first, m.First())
}

func TestTxSortedMapAdd(t *testing.T) {
	privateKey := types.GeneratePrivateKey()

	m := NewTxSortedMap()
	n := 100

	for i := 0; i < n; i++ {
		tx := types.NewRandomTransaction(privateKey)
		m.Add(tx)
		// cannot add the same twice
		m.Add(tx)

		assert.Equal(t, m.Count(), i+1)
		assert.True(t, m.Contains(tx.GetHash()))
		assert.Equal(t, len(m.lookup), m.txs.Len())
		assert.Equal(t, m.Get(tx.GetHash()), tx)
	}

	m.Clear()
	assert.Equal(t, m.Count(), 0)
	assert.Equal(t, len(m.lookup), 0)
	assert.Equal(t, m.txs.Len(), 0)
}

func TestTxSortedMapRemove(t *testing.T) {
	privateKey := types.GeneratePrivateKey()

	m := NewTxSortedMap()

	tx := types.NewRandomTransaction(privateKey)
	m.Add(tx)
	assert.Equal(t, m.Count(), 1)

	m.Remove(tx.GetHash())
	assert.Equal(t, m.Count(), 0)
	assert.False(t, m.Contains(tx.GetHash()))
}
