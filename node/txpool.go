package node

import (
	"fmt"
	"github.com/barreleye-labs/barreleye/common"
	"github.com/barreleye-labs/barreleye/core"
	"github.com/barreleye-labs/barreleye/core/types"
	"sync"
)

type TxPool struct {
	pending   *TxSortedMap
	maxLength int
}

func NewTxPool(maxLength int) *TxPool {
	return &TxPool{
		pending:   NewTxSortedMap(),
		maxLength: maxLength,
	}
}

func (p *TxPool) Add(tx *types.Transaction, chain *core.Blockchain) error {
	if p.pending.Contains(tx.GetHash()) {
		return common.ErrTransactionAlreadyPending
	}

	txs := p.Pending()
	for i := 0; i < len(txs); i++ {
		if txs[i].From == tx.From {
			txProcessed, err := chain.ReadTxByHash(txs[i].Hash)
			if err != nil {
				return err
			}

			if txProcessed == nil {
				return fmt.Errorf("this account already has a pending transaction")
			}

			p.pending.Remove(txs[i].Hash)
			break
		}
	}

	if p.pending.Count() == p.maxLength {
		oldest := p.pending.First()
		p.pending.Remove(oldest.GetHash())
	}
	p.pending.Add(tx)

	return nil
}

func (p *TxPool) Contains(hash common.Hash) bool {
	return p.pending.Contains(hash)
}

func (p *TxPool) Pending() []*types.Transaction {
	return p.pending.txs.Data
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

func (p *TxPool) PendingCount() int {
	return p.pending.Count()
}

type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[common.Hash]*types.Transaction
	txs    *common.List[*types.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[common.Hash]*types.Transaction),
		txs:    common.NewList[*types.Transaction](),
	}
}

func (t *TxSortedMap) First() *types.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	first := t.txs.Get(0)
	return t.lookup[first.GetHash()]
}

func (t *TxSortedMap) Get(h common.Hash) *types.Transaction {
	t.lock.RLock()
	defer t.lock.RUnlock()

	return t.lookup[h]
}

func (t *TxSortedMap) Add(tx *types.Transaction) {
	hash := tx.GetHash()

	t.lock.Lock()
	defer t.lock.Unlock()

	if _, ok := t.lookup[hash]; !ok {
		t.lookup[hash] = tx
		t.txs.Insert(tx)
	}
}

func (t *TxSortedMap) Remove(h common.Hash) {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.txs.Remove(t.lookup[h])
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

	t.lookup = make(map[common.Hash]*types.Transaction)
	t.txs.Clear()
}
