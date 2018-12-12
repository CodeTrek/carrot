package main

import (
	"fmt"

	"github.com/CodeTrek/carrot"
)

var target1 = func() int {
	return 1
}
var replacement1 = func() int {
	return 2
}
var original1 = func() (a int) { return }

func test1() {
	carrot.Patch(target1, replacement1, original1)
	target1()
	original1()

	carrot.UnpatchAll()

	target1()
}

var target2 = func(p1, p2 [2000]byte) int {
	fmt.Printf("target%s%s%s%s\n", p1[0:1], p2[0:1], p1[0:1], p2[0:1])
	return 1
}

var replacement2 = func(p1, p2 [2000]byte) int {
	fmt.Println("replacement")
	return 2
}

var original2 = func(p1, p2 [2000]byte) int { fmt.Println("original"); return 3 }

func test2() {

	fmt.Println("\nafter patch")
	var b = [2000]byte{0}
	carrot.Patch(target2, replacement2, original2)
	target2(b, b)
	replacement2(b, b)
	original2(b, b)

	carrot.UnpatchAll()

	fmt.Println("\nafter unpatch")
	target2(b, b)
	replacement2(b, b)
	original2(b, b)
}

var target3 = func() string {
	return "target"
}

var replacement3 = func() string {
	return "replacement"
}

var original3 = func() string {
	return "original"
}

func test3() {
	fmt.Println("\n[test3]")

	//	fmt.Println("target3")
	//	carrot.Disas(target3)

	//	carrot.Break()
	//	target3()
	carrot.Patch(target3, replacement3, original3)
	//	fmt.Println("target3")
	//	carrot.Disas(target3)
	//	carrot.Break()
	//	original3()

	//	fmt.Println("original")
	//	carrot.Disas(original3)

	//	return
	//	fmt.Println("1")

	fmt.Printf("after patch\ntarget=%s, replacement=%s, original=%s\n", target3(), replacement3(), original3())

	carrot.UnpatchAll()
	fmt.Printf("after unpatch\ntarget=%s, replacement=%s, original=%s\n", target3(), replacement3(), original3())

}

var t4 = func(p1 [2000]byte) int {
	fmt.Printf("%s%s%s%s", p1[0:1], p1[0:1], p1[0:1], p1[0:1])
	return 1
}

var r4 = func(p1 [2000]byte) int {
	return 2
}

var o4 = func(p1 [2000]byte) int { return 3 }

func test4() {
	var b = [2000]byte{0}
	fmt.Println("t4")
	carrot.Disas(t4)
	carrot.Patch(t4, r4, o4)
	fmt.Println("t4")
	carrot.Disas(t4)
	fmt.Println("r4")
	carrot.Disas(r4)
	//	fmt.Println("bridge")
	//	carrot.Disas(carrot.GetBridge(t4))
	carrot.Break()
	t4(b)
}

func main() {
	//	test1()
	//	test2()
	//	test3()
	test4()
}
