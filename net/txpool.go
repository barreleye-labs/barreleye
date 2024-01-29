package net

import (
	"github.com/barreleye-labs/barreleye/common"
	types2 "github.com/barreleye-labs/barreleye/core/types"
	"sync"
)

type TxPool struct {
	all     *TxSortedMap
	pending *TxSortedMap
	// 풀 사이즈
	// 풀이 가득차면 가장 오래된 트랜잭션부터 푸루닝할 것.
	maxLength int
}

func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

func (p *TxPool) Add(tx *types2.Transaction) {
	// 푸루닝
	if p.all.Count() == p.maxLength {
		oldest := p.all.First()
		p.all.Remove(oldest.GetHash(types2.TxHasher{}))
	}

	if !p.all.Contains(tx.GetHash(types2.TxHasher{})) {
		p.all.Add(tx)
		p.pending.Add(tx)
	}
}

func (p *TxPool) Contains(hash common.Hash) bool {
	return p.all.Contains(hash)
}

// Pending returns a slice of transactions that are in the pending pool
func (p *TxPool) Pending() []*types2.Transaction {
	return p.pending.txx.Data
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[common.Hash]*types2.Transaction
	txx    *common.List[*types2.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[common.Hash]*types2.Transaction),
		txx:    common.NewList[*types2.Transaction](),
	}
}

func (t *TxSortedMap) First() *types2.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.txx.Get(0)
	return t.lookup[first.GetHash(types2.TxHasher{})]
}

func (t *TxSortedMap) Get(h common.Hash) *types2.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

func (t *TxSortedMap) Add(tx *types2.Transaction) {
	hash := tx.GetHash(types2.TxHasher{})

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = tx
		t.txx.Insert(tx)
	}
}

func (t *TxSortedMap) Remove(h common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.txx.Remove(t.lookup[h])
	delete(t.lookup, h)
}

func (t *TxSortedMap) Count() int {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return len(t.lookup)
}

func (t *TxSortedMap) Contains(h common.Hash) bool {
	t.lock.RLock()
	defer t.lock.RUnlock()

	_, ok := t.lookup[h]
	return ok
}

func (t *TxSortedMap) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.lookup = make(map[common.Hash]*types2.Transaction)
	t.txx.Clear()
}
