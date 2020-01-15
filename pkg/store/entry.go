package store

import (
	"time"

	"cpl.li/go/errnil/pkg/errnil"
)

type Entry struct {
	UpdatedAt      time.Time         `json:"updated_at"`
	Repo           string            `json:"repo"`
	Positions      []errnil.Position `json:"positions"`
	PositionsCount int               `json:"positions_count"`
}
