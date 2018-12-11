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

var f2 = func(p1, p2 [2000]byte) int {
	fmt.Printf("original%s%s\n", p1[0:1], p2[0:1])
	return 1
}

var newF2 = func(p1, p2 [2000]byte) int {
	fmt.Println("patched")
	return 2
}

var oldF2 = func(p1, p2 [2000]byte) int { return 3 }

func main() {
	//	carrot.Patch(f1, newF1, ff)
	//	fmt.Printf("f1=%d, newF1=%d, ff=%d\n", f1(), newF1(), ff())

	fmt.Println("\nnewF2")
	carrot.Disas(newF2)

	fmt.Println("\nf2")
	carrot.Disas(f2)
	fmt.Println("\noldF2")
	carrot.Disas(oldF2)
	var b = [2000]byte{0}
	carrot.Patch(f2, newF2, oldF2)
	fmt.Println("\nf2")
	carrot.Disas(f2)
	fmt.Println("\noldF2")
	carrot.Disas(oldF2)

	fmt.Printf("f2=%d\n", f2(b, b))
	fmt.Printf("newF2=%d\n", newF2(b, b))
	fmt.Printf("oldF2=%d\n", oldF2(b, b))
	return
}
