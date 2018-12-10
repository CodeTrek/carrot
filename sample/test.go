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

func f2(p1, p2, p3, p4 [2000]byte) int {
	fmt.Printf("original%s%s%s%s\n", p1[0:1], p3[0:1], p3[0:1], p4[0:1])
	return 1
}

func newF2(p1, p2, p3, p4 [2000]byte) int {
	fmt.Println("patched")
	return 2
}

var oldF2 = func(p1, p2, p3, p4 [2000]byte) int { return 3 }

func main() {
	carrot.Patch(f1, newF1, ff)

	fmt.Printf("%d", f1())

	return
	//	carrot.Patch(f2, newF2, oldF2)
	//	var p = [2000]byte{0}
	//	f2(p, p, p, p)
}
