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
	fmt.Printf("original%s%s%s%s\n", p1[0:1], p2[0:1], p1[0:1], p2[0:1])
	return 1
}

var newF2 = func(p1, p2 [2000]byte) int {
	fmt.Println("patched")
	return 2
}

var oldF2 = func(p1, p2 [2000]byte) int { fmt.Println("old"); return 3 }

func testF1() {
	fmt.Println("\nf1")
	carrot.Disas(f1)

	carrot.Patch(f1, newF1, ff)
	f1()
	ff()

	fmt.Println("\nf1")
	carrot.Disas(f1)

	carrot.UnpatchAll()

	fmt.Println("\nf1")
	carrot.Disas(f1)

	f1()
}

func testF2() {
	//	fmt.Println("\nf2")
	//	carrot.Disas(f2)

	var b = [2000]byte{0}
	carrot.Patch(f2, newF2, oldF2)
	f2(b, b)
	oldF2(b, b)
	newF2(b, b)

	//	fmt.Println("\nf2")
	//	carrot.Disas(f2)

	carrot.UnpatchAll()

	//	fmt.Println("\nf2")
	//	carrot.Disas(f2)

	f2(b, b)
	newF2(b, b)
	oldF2(b, b)
}

func main() {
	testF2()
}
