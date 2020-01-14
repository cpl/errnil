package errnil

import (
	"fmt"
	"go/token"
)

// Position is a custom representation of the token.Position struct. It allows for custom methods and other additions.
type Position struct {
	Filename string `json:"filename"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Offset   int    `json:"offset"`
}

func positionFromTokenPosition(position token.Position) Position {
	return Position{
		Filename: position.Filename,
		Line:     position.Line,
		Column:   position.Column,
		Offset:   position.Offset,
	}
}

func (p Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.Filename, p.Line, p.Column)
}
