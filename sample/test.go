package main

import (
	"fmt"

	"github.com/CodeTrek/carrot"
)

func f() {
}
func newF() {
}

var f0 = func() {}

var f1 = func() int {
	return 1
}
var newF1 = func() int {
	return 2
}

var ff = func() (a int) { return }

func f2(p1 [2000]byte) int {
	fmt.Printf("original%s\n", p1[0:1])
	return 1
}

func newF2(p1 [2000]byte) int {
	fmt.Println("patched")
	return 2
}

var oldF2 = func(p1 [2000]byte) int { return 3 }

func main() {
	//	carrot.Patch(f1, newF1, ff)
	//	fmt.Printf("f1=%d, newF1=%d, ff=%d\n", f1(), newF1(), ff())

	var b = [2000]byte{0}
	carrot.Patch(f2, newF2, oldF2)
	fmt.Printf("f2=%d, newF2=%d, oldF2=%d\n", f2(b), newF2(b), oldF2(b))
	return
}
