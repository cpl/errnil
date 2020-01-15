package store

import (
	"errors"

	"cpl.li/go/errnil/pkg/errnil"
)

type Store interface {
	SetPositions(repo string, positions []errnil.Position) error
	GetPositions(repo string) ([]errnil.Position, error)
	GetPositionsCount(repo string) (int, error)
}

var (
	ErrRepoNotFound = errors.New("repo not found")
)
