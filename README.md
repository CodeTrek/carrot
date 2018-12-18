# GO carrot patch

This is a detours like library for Golang!

## Ref

* https://github.com/vmt/udis86
* https://github.com/bouk/monkey

## Design

```
  call target                                                  call original
      v                                                              |
      |                                                              |
      |       target                              replacement        |       original
      +-> +-----------+                     +-> +-------------+      +---> +----------+
          |    JMP1   | jump to replacement |   |     new     |            |    JMP   |
      +-> +-----------+---------------------+   |    target   |            +----+-----+
      |   |    code   |       (jmp2r)           |     func    |                 |
      |   |    ...    |                         |             |                 |
      |   +-----------+                         +-------------+          jump to bridge
      |   |    JMP2   |                                                         |
      |   +-----------+---------+                                               |
    (jmp2t)                     |                                             (jmp2b)
      |                  jump to incr stack                                     |
      |                       (jmp2i)                                           |
      |          bridge         |                                               |
      |   +----------------+ <--------------+-----------------------------------+
      |   | code from JMP1 |    |           ^
      |   +----------------+    |           |
      |   |    JMP         |    |           |
      +---+----------------+    |  continue exec after stack incrd
          +----------------+ <--+         (jmp2b)
          | call incr stk  |                |
          +----------------+                |
          |       JMP      |                |
          +----------------+----------------+

```

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
