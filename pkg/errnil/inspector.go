package errnil

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	stdpath "path"
	"path/filepath"
)

func Inspect(path string) ([]token.Position, error) {
	var positions []token.Position

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

	return positions, nil
}

func extractPositions(tokenSet *token.FileSet, astFp *ast.File) (positions []token.Position) {
	ast.Inspect(astFp, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.BinaryExpr:
			if n.Op != token.NEQ {
				break
			}
			x, ok := n.X.(*ast.Ident)
			if !ok {
				break
			}
			y, ok := n.Y.(*ast.Ident)
			if !ok {
				break
			}

			if (x.Name == "err" && y.Name == "nil") || (x.Name == "nil" && y.Name == "err") {
				pos := tokenSet.PositionFor(n.Pos(), true)
				positions = append(positions, pos)
			}
		}

		return true
	})

	return
}
