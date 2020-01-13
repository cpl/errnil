package main

import (
	"fmt"
	"runtime"

	"cpl.li/go/errnil/pkg/errnil"
)

func main() {
	counter := errnil.NewCounter(runtime.NumCPU())

	// testing with downloading packages in order to offer this as a service
	//dir, err := errnil.Download("go.uber.org/zap", path.Join(os.TempDir(), "errnil"))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(dir)

	count, err := counter.Count(".")
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
}
