package main

import (
	"fmt"
	"path"
	"strings"

	"cpl.li/go/errnil/pkg/store"

	"cpl.li/go/errnil/pkg/errnil"
)

const (
	shieldsEndpoint = "https://img.shields.io/static/v1"
	badgeColor      = "e44"
	badgeLabel      = "err != nil"
)

func fmtBadgeURL(message, color, style string) string {
	return fmt.Sprintf("%s?label=%s&message=%s&color=%s&style=%s",
		shieldsEndpoint, badgeLabel, message, color, style)
}

func countSourceFiles(positions []errnil.Position) map[string]int {
	files := make(map[string]int)

	for _, position := range positions {
		if _, ok := files[position.Filename]; !ok {
			files[position.Filename] = 1
		} else {
			files[position.Filename] += 1
		}
	}

	return files
}

func cleanPositions(positions []errnil.Position, tmpDir string) {
	for idx := range positions {
		positions[idx].Filename = strings.Replace(positions[idx].Filename, tmpDir+"/", "", 1)
	}
}

func downloadInspectStore(repo, downloadDir string, storage store.Store) (store.Entry, error) {
	var entry store.Entry

	downloadPath, err := errnil.Download(repo, downloadDir)
	if err != nil {
		return entry, fmt.Errorf("failed downloading repo, %s", err.Error())
	}

	positions, err := errnil.Inspect(downloadPath)
	if err != nil {
		return entry, fmt.Errorf("failed inspecting repo, %s", err.Error())
	}

	cleanPositions(positions, path.Join(downloadPath, ".."))

	entry, err = storage.SetEntry(repo, positions)
	if err != nil {
		return entry, fmt.Errorf("failed storing entry, %s", err.Error())
	}

	return entry, nil
}
