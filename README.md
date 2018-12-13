# GO carrot patch

This is a detours like library for Golang!

## Ref

https://github.com/vmt/udis86
https://github.com/bouk/monkey

## Using carrot

```go
package main

import (
	"fmt"

	"github.com/CodeTrek/carrot"
)

var t3 = func() string { return "target" }
var r3 = func() string { return "replacement" }
var o3 = func() string { return "original" }

func main() {
	carrot.Patch(t3, r3, o3)
	fmt.Printf("after patch\ntarget=%s, replacement=%s, original=%s\n", t3(), r3(), o3())
    
    carrot.Unpatch(t3)
}

```

## Warning!!

This library is not fully tested.
