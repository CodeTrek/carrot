package carrot

/*
#include "thirdparty/udis86/libudis86/decode.c"
#include "thirdparty/udis86/libudis86/itab.c"
#include "thirdparty/udis86/libudis86/udis86.c"
#include "thirdparty/udis86/libudis86/syn.c"
#include "thirdparty/udis86/libudis86/syn-intel.c"
//#include "thirdparty/udis86/libudis86/syn-att.c"
#include "udis_code.c"
*/
import "C"
import (
	"reflect"
	"unsafe"
)

func memoryAccess(p uintptr, len int) []byte {
	// 这里直接给足够大的访问权限，只要不是真的去读取就不会有问题
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  len,
		Cap:  len,
	}))
}

func disas(bytes []byte) {
	var code = (*C.uint8_t)(unsafe.Pointer(&bytes[0]))
	var length = (C.size_t)(len(bytes))

	u := C.udis_init(code, length)
	C.udis_print(u)
	C.udis_final(u)
}
