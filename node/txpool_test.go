package node

import (
	"github.com/barreleye-labs/barreleye/core/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxMaxLength(t *testing.T) {
	privateKey := types.GeneratePrivateKey()
	p := NewTxPool(1)
	p.Add(types.NewRandomTransaction(privateKey))
	assert.Equal(t, 1, p.all.Count())

	p.Add(types.NewRandomTransaction(privateKey))
	p.Add(types.NewRandomTransaction(privateKey))
	p.Add(types.NewRandomTransaction(privateKey))
	tx := types.NewRandomTransaction(privateKey)
	p.Add(tx)
	assert.Equal(t, 1, p.all.Count())
	assert.True(t, p.Contains(tx.GetHash()))
}

func TestTxPoolAdd(t *testing.T) {
	privateKey := types.GeneratePrivateKey()

	p := NewTxPool(11)
	n := 10

	for i := 1; i <= n; i++ {
		tx := types.NewRandomTransaction(privateKey)
		p.Add(tx)
		// cannot add twice
		p.Add(tx)

		assert.Equal(t, i, p.PendingCount())
		assert.Equal(t, i, p.pending.Count())
		assert.Equal(t, i, p.all.Count())
	}
}

func TestTxPoolMaxLength(t *testing.T) {
	privateKey := types.GeneratePrivateKey()

	maxLen := 10
	p := NewTxPool(maxLen)
	n := 100
	txx := []*types.Transaction{}

	for i := 0; i < n; i++ {
		tx := types.NewRandomTransaction(privateKey)
		p.Add(tx)

		if i > n-(maxLen+1) {
			txx = append(txx, tx)
		}
	}

	assert.Equal(t, p.all.Count(), maxLen)
	assert.Equal(t, len(txx), maxLen)

	for _, tx := range txx {
		assert.True(t, p.Contains(tx.GetHash()))
	}
}

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
		assert.Equal(t, len(m.lookup), m.txx.Len())
		assert.Equal(t, m.Get(tx.GetHash()), tx)
	}

	m.Clear()
	assert.Equal(t, m.Count(), 0)
	assert.Equal(t, len(m.lookup), 0)
	assert.Equal(t, m.txx.Len(), 0)
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
