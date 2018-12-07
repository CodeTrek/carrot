package main

import (
	"fmt"

	"github.com/CodeTrek/carrot"
)

func f1() int {
	fmt.Println("hello")
	return 1
}

func newF1() int {
	return 2
}

var ff = func() (a int) { return }

func main() {
	carrot.Patch(f1, newF1, ff)
}
