# carrot

 call target                                               call original
     v                                                           v
     |                                                           |
     |       target                           replacement        |       original
     +-> +-----------+                  +-> +-------------+      +---> +----------+
 . . . . |    JMP    | jump to replaced |   |     new     |            |    JMP   |
 .   +-> +-----------+------------------+   |    target   |            +----+-----+
 .   |   |    code   |                      |     func    |                 |
 .   |   |    ...    |                      |             |                 |
 .   |   +-----------+                      +-------------+          jump to bridge
 .   |                                                                      |
 .   |                                                                      |
 .   |      bridge                                                          |
 .   |   +-----------+ <----------------------------------------------------+
 .   |   |    code   | < . . . .
 .   |   +-----------+         .
 .   |   |    JMP    |         .
 .   +---+-----------+         .
 . . . . . . . . . . . . . . . .
 copied from the beginning of target
  (instruction replaced by JMP)
