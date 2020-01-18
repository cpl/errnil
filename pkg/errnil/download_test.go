package errnil

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	t.Parallel()

	tmpDir := path.Join(os.TempDir(), "errnil_download_test")
	defer os.RemoveAll(tmpDir)

	_, err := Download("https://does.not.exist.local", tmpDir)
	assert.EqualError(t, err, "failed finding repo, invalid import path \"https://does.not.exist.local\"")

	downloadPath, err := Download("cpl.li/go/alpha", tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, path.Join(tmpDir, "cpl.li", "go", "alpha"), downloadPath)

	downloadPath, err = Download("cpl.li/go/alpha", tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, path.Join(tmpDir, "cpl.li", "go", "alpha"), downloadPath)

	assert.NoError(t, os.RemoveAll(downloadPath))

	downloadPath, err = Download("cpl.li/go/alpha", tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, path.Join(tmpDir, "cpl.li", "go", "alpha"), downloadPath)
}
