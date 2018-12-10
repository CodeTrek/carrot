# carrot

 call target                                                  call original
     v                                                              |
     |                                                              |
     |       target                              replacement        v       original
     +-> +-----------+                     +-> +-------------+      +---> +----------+
 . . . . |    JMP    | jump to replacement |   |     new     |      ^     |    JMP   |
 .   +-> +-----------+---------------------+   |    target   |      |     +----+-----+
 .   |   |    code   |                         |     func    |      |          |
 .   |   |    ...    |                         |             |      |          |
 .   |   +-----------+                         +-------------+      |   jump to bridge
 .   |   | more stack|      jump after adjust stack                 |          |
 .   |   +-----------+----------------------------------------------+          |
 .   |      bridge                                                             |
 .   |   +-----------+ <-------------------------------------------------------+
 .   |   |    code   | < . . . .
 .   |   +-----------+         .
 .   |   |    JMP    |         .
 .   +---+-----------+         .
 . . . . . . . . . . . . . . . .
 copied from the beginning of target
  (instruction replaced by JMP)
