package carrot

/*
#cgo CFLAGS: -DBRIDGE_PIECE_SIZE=32
#cgo CFLAGS: -DBRIDGE_PAGE_COUNT=5
#include "thirdparty/udis86/libudis86/decode.c"
#include "thirdparty/udis86/libudis86/itab.c"
#include "thirdparty/udis86/libudis86/udis86.c"
#include "thirdparty/udis86/libudis86/syn.c"
#include "thirdparty/udis86/libudis86/syn-intel.c"
//#include "thirdparty/udis86/libudis86/syn-att.c"
#include "udis_code.c"

inline void int3()
{
	__asm__("int $0x03");
}

*/
import "C"
import "unsafe"

func udisDisas(data []byte) {
	var code = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	var length = (C.size_t)(len(data))
	C.udis_disas(code, length)
}

func backupInstruction(location uintptr, minLen int) (code []byte, moreStackPtr uintptr, funcEnded bool) {
	v := C.udis_backup_instruction(*(**C.uint8_t)(unsafe.Pointer(&location)), 300, C.size_t(minLen))
	if v.success == 0 {
		udisDisas(memoryAccess(location, 300))
		panic("backup instruction failed!")
	}
	code = memoryAccess(location, int(v.code_len))
	moreStackPtr = uintptr(v.adjust_stack_jmp)
	funcEnded = v.reach_end == 1
	return
}

// Break int3
func Break() {
	C.int3()
}
