package main

import (
	"fmt"

	"../../carrot"
)

var t1 = func() int { return 1 }
var r1 = func() int { return 2 }
var o1 = func() (a int) { return }

var t2 = func(p1, p2 [2000]byte) int {
	fmt.Printf("target%s%s%s%s\n", p1[0:1], p2[0:1], p1[0:1], p2[0:1])
	return 1
}
var r2 = func(p1, p2 [2000]byte) int { fmt.Println("replacement"); return 2 }
var o2 = func(p1, p2 [2000]byte) int { fmt.Println("original"); return 3 }

var t3 = func() string { return "target" }
var r3 = func() string { return "replacement" }
var o3 = func() string { return "original" }

var t4 = func(p1 [2000]byte) int { fmt.Printf("%s%s%s%s", p1[0:1], p1[0:1], p1[0:1], p1[0:1]); return 1 }
var r4 = func(p1 [2000]byte) int { return 2 }
var o4 = func(p1 [2000]byte) int { return 3 }

func test1() {
	carrot.Patch(t1, r1, o1)
	t1()
	o1()

	carrot.UnpatchAll()

	t1()
}

func test2() {

	fmt.Println("\nafter patch")
	var b = [2000]byte{0}
	carrot.Patch(t2, r2, o2)
	t2(b, b)
	r2(b, b)
	o2(b, b)

	carrot.UnpatchAll()

	fmt.Println("\nafter unpatch")
	t2(b, b)
	r2(b, b)
	o2(b, b)
}

func test3() {
	fmt.Println("\n[test3]")

	fmt.Printf("before patch\ntarget=%s, replacement=%s, original=%s\n", t3(), r3(), o3())

	carrot.Patch(t3, r3, o3)
	fmt.Printf("after patch\ntarget=%s, replacement=%s, original=%s\n", t3(), r3(), o3())

	carrot.UnpatchAll()
	fmt.Printf("after unpatch\ntarget=%s, replacement=%s, original=%s\n", t3(), r3(), o3())

}

func test4() {
	fmt.Println("\n[test4]")
	var b = [2000]byte{0}
	fmt.Printf("before patch\ntarget=%d, replacement=%d, original=%d\n", t4(b), r4(b), o4(b))

	carrot.Patch(t4, r4, o4)
	fmt.Printf("after patch\ntarget=%d, replacement=%d, original=%d\n", t4(b), r4(b), o4(b))

	carrot.UnpatchAll()
	fmt.Printf("after unpatch\ntarget=%d, replacement=%d, original=%d\n", t4(b), r4(b), o4(b))
}

func test5() {
	fmt.Println("\n[test5]")
	carrot.Patch(t3, r3, o3)
	//	carrot.Break()
	fmt.Println(carrot.Patch(t4, r4, o4))
}

func main() {
	test1()
	test2()
	test3()
	test4()
	test5()
}
