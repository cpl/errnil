package badger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ToMarkdown(t *testing.T) {
	t.Parallel()

	markdown := ToMarkdown("cpl.li/go/errnil", "")
	assert.Equal(t, "[![err != nil](https://img.shields.io/endpoint?url=https%3A%2F%2Ferrnil.cpl.li%2Fapi%2Fbadge%3Frepo%3Dcpl.li%252Fgo%252Ferrnil%26style%3Dflat)](https://errnil.cpl.li/api/inspect?repo=cpl.li/go/errnil)", markdown)
}
