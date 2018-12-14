package carrot

/*
#include "./thirdparty/udis86/libudis86/decode.c"
#include "./thirdparty/udis86/libudis86/itab.c"
#include "./thirdparty/udis86/libudis86/udis86.c"
#include "./thirdparty/udis86/libudis86/syn.c"
#include "./thirdparty/udis86/libudis86/syn-att.c"
#include "./udis_copy_instruction.c"
#include "./udis_code.c"
*/
import "C"
import (
	"unsafe"
)

func udisDisas(data []byte) {
	var code = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	var length = (C.size_t)(len(data))
	C.udis_disas(code, length)
}

func backupInstruction(src uintptr, jmpLen int) (code []byte, incrStk []byte, srcCopiedLen int, incrStkPtr uintptr, funcEnded bool) {
	v := C.udis_backup_instruction(*(**C.uint8_t)(unsafe.Pointer(&src)), 300, C.size_t(jmpLen))
	if v.success == 0 {
		udisDisas(memoryAccess(src, 300))
		panic("backup instruction failed!")
	}
	code = make([]byte, int(v.copied_len))
	incrStk = make([]byte, int(v.incr_stack_len))
	copy(code, memoryAccess(uintptr(unsafe.Pointer(&v.copied[0])), int(v.copied_len)))
	copy(incrStk, memoryAccess(uintptr(unsafe.Pointer(&v.incr_stack[0])), int(v.incr_stack_len)))
	incrStkPtr = uintptr(v.incr_stack_ptr)
	srcCopiedLen = int(v.copied_src_len)
	funcEnded = v.reach_end == 1
	return
}
