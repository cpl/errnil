package main

import (
	"fmt"
	"os"

	"cpl.li/go/errnil/pkg/errnil"
)

func main() {
	positions, err := errnil.Inspect(os.Args[1])
	if err != nil {
		panic(err)
	}

	for _, pos := range positions {
		fmt.Printf("%s:%d:%d\n", pos.Filename, pos.Line, pos.Column)
	}
}
