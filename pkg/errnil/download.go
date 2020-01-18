package errnil

import (
	"fmt"
	"os"
	"path"

	"golang.org/x/tools/go/vcs"
)

// Download will take the given package repository and download it to a destination. The returned value is the final
// path containing the package and including the parent directories. If the package already exists on disk, then it will
// be updated.
func Download(repo, destination string) (downloadPath string, err error) {
	return download(repo, destination)
}

func download(repo, destination string) (string, error) {
	repoRoot, err := vcs.RepoRootForImportPath(repo, false)
	if err != nil {
		return "", fmt.Errorf("failed finding repo, %w", err)
	}

	dir := path.Join(destination, repoRoot.Root)

	if exists(dir) {
		if err := repoRoot.VCS.Download(dir); err != nil {
			return "", fmt.Errorf("failed updating repo, %w", err)
		}
	} else {
		if err := os.MkdirAll(dir, 0766); err != nil {
			return "", fmt.Errorf("failed making path, %w", err)
		}

		if err := repoRoot.VCS.Create(dir, repoRoot.Repo); err != nil {
			return "", fmt.Errorf("failed obtaining repo, %w", err)
		}
	}

	return dir, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return !os.IsNotExist(err)
}
