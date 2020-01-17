package store

import (
	"testing"

	"cpl.li/go/errnil/pkg/errnil"

	"github.com/stretchr/testify/assert"
)

func TestPrimitiveStore(t *testing.T) {
	t.Parallel()

	var entry Entry
	var err error

	repo := "https://github.com/example/example"
	testPosition := errnil.Position{
		Filename: "example/pkg/test.go",
		Line:     10,
		Column:   20,
		Offset:   300,
	}

	storage := NewPrimitiveStore()
	entry, err = storage.GetEntry(repo)
	assert.EqualError(t, err, ErrRepoNotFound.Error())

	entry, err = storage.SetEntry(repo, []errnil.Position{
		testPosition,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, entry.PositionsCount)
	assert.Equal(t, repo, entry.Repo)

	storeEntry, err := storage.GetEntry(repo)
	assert.NoError(t, err)
	assert.Equal(t, entry, storeEntry)

	entry, err = storage.SetEntry(repo, []errnil.Position{testPosition, testPosition, testPosition})
	assert.NoError(t, err)
	assert.Equal(t, 3, entry.PositionsCount)
	assert.Equal(t, repo, entry.Repo)
	assert.Equal(t, testPosition, entry.Positions[0])
}
