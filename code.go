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
	"container/list"
	"reflect"
	"syscall"
	"unsafe"
)

const piceSize int = 32

var (
	freeBridgeList = initBridge(4)
	usedBridgeMap  = make(map[uintptr][]byte)
)

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

func initBridge(pageCount int) *list.List {
	ptr := uintptr(C.udis_init_and_get_bridge(C.int(pageCount), C.int(syscall.Getpagesize())))
	if ptr == 0 {
		panic("init failed\n")
	}

	size := pageCount * syscall.Getpagesize()

	l := list.New()
	for offset := 0; offset < size; offset += piceSize {
		var mem uintptr = ptr + uintptr(offset)
		l.PushBack(memoryAccess(mem, piceSize))
	}

	return l
}

type value struct {
	_   uintptr
	ptr unsafe.Pointer
}

func location(v reflect.Value) uintptr {
	return uintptr((*value)(unsafe.Pointer(&v)).ptr)
}

func memoryAccess(p uintptr, len int) []byte {
	// 这里直接给足够大的访问权限，只要不是真的去读取就不会有问题
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  len,
		Cap:  len,
	}))
}

func allocBridge() []byte {
	if freeBridgeList.Len() <= 0 {
		panic("no more pices!")
	}

	node := freeBridgeList.Front()
	v := node.Value
	freeBridgeList.Remove(node)

	vv, ok := v.([]byte)
	if !ok {
		panic("type failed")
	}

	ptr := uintptr(unsafe.Pointer(&vv[0]))
	usedBridgeMap[ptr] = vv

	return vv
}

func freeBridge(pice []byte) {
	ptr := uintptr(unsafe.Pointer(&pice[0]))
	vv, ok := usedBridgeMap[ptr]
	if !ok {
		panic("ptr not found!")
	}

	delete(usedBridgeMap, ptr)
	freeBridgeList.PushBack(vv)
}

func disas(bytes []byte) {
	var code = (*C.uint8_t)(unsafe.Pointer(&bytes[0]))
	var length = (C.size_t)(len(bytes))

	u := C.udis_init(code, length)
	C.udis_print(u)
	C.udis_final(u)
}
