package carrot

import (
	"bytes"
	"container/list"
	"reflect"
	"syscall"
	"unsafe"
)

const bridgePiceSize = 32

var (
	freeBridgeList = bridge()
	usedBridgeMap  = make(map[uintptr][]byte)
)

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

func bridge() *list.List {
	ls := list.New()
	allocBridge(ls)
	return ls
}

func allocBridge(ls *list.List) {
	data := allocPage()

	g := makeWritable(uintptr(unsafe.Pointer(&data[0])), len(data))
	defer g.restore()

	// split whole bridge buffer to pieces
	size := len(data)
	p := bytes.Repeat([]byte{0xcc}, bridgePiceSize)
	for i := 0; i < size; i += bridgePiceSize {
		var b = data[i : i+bridgePiceSize]
		copy(b, p)
		ls.PushBack(b)
	}
}

type value struct {
	_   uintptr
	ptr unsafe.Pointer
}

func locationFunc(v reflect.Value) uintptr {
	return uintptr((*value)(unsafe.Pointer(&v)).ptr)
}

type vs struct {
	ptr unsafe.Pointer
}

func locationSlice(s *[]byte) uintptr {
	return uintptr((*vs)(unsafe.Pointer(s)).ptr)
}

func memoryAccess(p uintptr, len int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: p,
		Len:  len,
		Cap:  len,
	}))
}

func allocBridgePiece() []byte {
	if freeBridgeList.Len() <= 0 {
		allocBridge(freeBridgeList)
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

func freeBridgePiece(piece []byte) {
	ptr := uintptr(unsafe.Pointer(&piece[0]))
	vv, ok := usedBridgeMap[ptr]
	if !ok {
		panic("ptr not found!")
	}

	delete(usedBridgeMap, ptr)
	freeBridgeList.PushBack(vv)
}
