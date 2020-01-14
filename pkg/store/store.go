package store

import "cpl.li/go/errnil/pkg/errnil"

type Store interface {
	SetPositions(repo string, positions []errnil.Position) error
	GetPositions(repo string) ([]errnil.Position, error)
	GetPositionsCount(repo string) (int, error)
}
