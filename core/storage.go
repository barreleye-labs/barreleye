package core

import "github.com/barreleye-labs/barreleye/core/types"

type Storage interface {
	Put(*types.Block) error
}

type MemoryStore struct {
}

func NewMemorystore() *MemoryStore {
	return &MemoryStore{}
}

func (s *MemoryStore) Put(b *types.Block) error {
	return nil
}
