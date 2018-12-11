# carrot

 call target                                                  call original
     v                                                              |
     |                                                              |
     |       target                              replacement        |       original
     +-> +-----------+                     +-> +-------------+      +---> +----------+
 . . . . |    JMP    | jump to replacement |   |     new     |            |    JMP   |
 .   +-> +-----------+---------------------+   |    target   |            +----+-----+
 .   |   |    code   |       (jmp2r)           |     func    |                 |
 .   |   |    ...    |                         |             |                 |
 .   |   +-----------+                         +-------------+          jump to bridge
 .   |   | more stack|                                                         |
 .   |   +-----------+------+                                                  |
 . (jmp2t)                  |  jump after adjust stack                       (jmp2b)
 .   |                   (jmp2b)                                               |
 .   |                      |                                                  |
 .   |      bridge          v                                                  |
 .   |   +-----------+ <----+--------------------------------------------------+
 .   |   |    code   | < . . . .
 .   |   +-----------+         .
 .   |   |    JMP    |         .
 .   +---+-----------+         .
 . . . . . . . . . . . . . . . .
 copied from the beginning of target
  (instruction replaced by JMP)
