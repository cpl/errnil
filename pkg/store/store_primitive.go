package store

import (
	"sync"

	"cpl.li/go/errnil/pkg/errnil"
)

type PrimitiveStore struct {
	dataLock sync.RWMutex
	data     map[string][]errnil.Position
}

func NewPrimitiveStore() Store {
	return &PrimitiveStore{
		data: make(map[string][]errnil.Position),
	}
}

func (p PrimitiveStore) SetPositions(repo string, positions []errnil.Position) error {
	p.dataLock.Lock()
	defer p.dataLock.Unlock()
	p.data[repo] = positions
	return nil
}

func (p PrimitiveStore) GetPositions(repo string) ([]errnil.Position, error) {
	p.dataLock.RLock()
	defer p.dataLock.RUnlock()

	positions, ok := p.data[repo]
	if !ok {
		return nil, ErrRepoNotFound
	}

	return positions, nil
}

func (p PrimitiveStore) GetPositionsCount(repo string) (int, error) {
	p.dataLock.RLock()
	defer p.dataLock.RUnlock()

	positions, ok := p.data[repo]
	if !ok {
		return -1, ErrRepoNotFound
	}

	return len(positions), nil
}
