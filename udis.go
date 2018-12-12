package carrot

/*
#include "thirdparty/udis86/libudis86/decode.c"
#include "thirdparty/udis86/libudis86/itab.c"
#include "thirdparty/udis86/libudis86/udis86.c"
#include "thirdparty/udis86/libudis86/syn.c"
#include "thirdparty/udis86/libudis86/syn-att.c"
#include "udis_copy_instruction.c"
#include "udis_code.c"
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

func backupInstruction(location uintptr, minLen int, dataPtr uintptr) (code []byte, dataUsed []byte, srcCopiedLen int, moreStackPtr uintptr, funcEnded bool) {
	v := C.udis_backup_instruction(*(**C.uint8_t)(unsafe.Pointer(&location)), 300, C.size_t(minLen), C.uintptr_t(dataPtr))
	if v.success == 0 {
		udisDisas(memoryAccess(location, 300))
		panic("backup instruction failed!")
	}
	code = make([]byte, int(v.copied_len))
	dataUsed = make([]byte, int(v.data_len))
	copy(code, memoryAccess(uintptr(unsafe.Pointer(&v.copied[0])), int(v.copied_len)))
	copy(dataUsed, memoryAccess(uintptr(unsafe.Pointer(&v.data[0])), int(v.data_len)))
	moreStackPtr = uintptr(v.adjust_stack_jmp)
	srcCopiedLen = int(v.copied_src_len)
	funcEnded = v.reach_end == 1
	return
}
