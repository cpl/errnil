package store

import "cpl.li/go/errnil/pkg/errnil"

type Store interface {
	SetEntry(repo string, positions []errnil.Position) (Entry, error)
	GetEntry(repo string) (Entry, error)
}
