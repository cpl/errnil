package errnil

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInspector(t *testing.T) {
	t.Parallel()

	tmpDir := path.Join(os.TempDir(), "errnil_test")
	assert.NoError(t, os.MkdirAll(tmpDir, 0766))
	defer os.RemoveAll(tmpDir)

	fp, err := os.Create(path.Join(tmpDir, "errnil.go"))
	assert.NoError(t, err)
	defer fp.Close()

	_, err = fp.WriteString(`package main

	import "fmt"

	func main() {
		fmt.Println("hello world")
		err := fmt.Errorf("testing testing")
		if err != nil {
			fmt.Println(err)
		}
		if err != nil {
			fmt.Println(err)
		}

		if err := filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
			if f.IsDir() || stdpath.Ext(f.Name()) != ".go" {
				return nil
			}
	
			fset := token.NewFileSet()
	
			fast, err := parser.ParseFile(fset, p, nil, parser.DeclarationErrors)
			if err != nil {
				return fmt.Errorf("failed parsing file, %w", err)
			}
	
			positions = append(positions, extractPositions(fset, fast)...)
	
			return nil
	
		}); err != nil {
			return nil, fmt.Errorf("failed path traversal, %w", err)
		}
	}`)
	assert.NoError(t, err)

	positions, err := Inspect(tmpDir)
	assert.NoError(t, err)
	assert.Len(t, positions, 4)

	assert.Equal(t, path.Join(tmpDir, "errnil.go")+":8:6", positions[0].String())
}
