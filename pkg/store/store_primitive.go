package store

import (
	"sync"
	"time"

	"cpl.li/go/errnil/pkg/errnil"
)

type PrimitiveStore struct {
	dataLock *sync.RWMutex
	data     map[string]Entry
}

func NewPrimitiveStore() Store {
	return &PrimitiveStore{
		dataLock: new(sync.RWMutex),
		data:     make(map[string]Entry),
	}
}

func (p PrimitiveStore) SetEntry(repo string, positions []errnil.Position) (Entry, error) {
	p.dataLock.Lock()
	defer p.dataLock.Unlock()

	entry := Entry{
		UpdatedAt:      time.Now().UTC(),
		Repo:           repo,
		Positions:      positions,
		PositionsCount: len(positions),
	}
	p.data[repo] = entry

	return entry, nil
}

func (p PrimitiveStore) GetEntry(repo string) (Entry, error) {
	p.dataLock.RLock()
	defer p.dataLock.RUnlock()

	entry, ok := p.data[repo]
	if !ok {
		return Entry{}, ErrRepoNotFound
	}

	return entry, nil
}
