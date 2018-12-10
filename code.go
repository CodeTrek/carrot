package carrot

import (
	"container/list"
	"reflect"
	"syscall"
	"unsafe"
)

var (
	freeBridgeList = initBridge()
	usedBridgeMap  = make(map[uintptr][]byte)
)

func pageStart(ptr uintptr) uintptr {
	return ptr & ^(uintptr(syscall.Getpagesize() - 1))
}

func initBridge() *list.List {
	bridgePtr, bridgeLen := udisInitBridge()
	pageSize := uintptr(syscall.Getpagesize())
	if bridgePtr == 0 {
		panic("init failed\n")
	}

	// align to pagesize
	bridgeStart := pageStart(bridgePtr)
	if bridgeStart < bridgePtr {
		bridgeStart += pageSize
	}

	// change page protection to executable
	pageCount := uintptr((bridgePtr + bridgeLen - bridgeStart) / pageSize)
	for i := uintptr(0); i < pageCount; i++ {
		makePageExecutable(bridgeStart + i*pageSize)
	}

	// split whole bridge buffer to pieces
	ls := list.New()
	size := pageCount * pageSize
	pieceSize := uintptr(udisBridgePieceSize())
	for offset := uintptr(0); offset < size; offset += pieceSize {
		var ptr = bridgeStart + offset
		ls.PushBack(memoryAccess(ptr, int(pieceSize)))
	}

	return ls
}

type value struct {
	_   uintptr
	ptr unsafe.Pointer
}

func location(v reflect.Value) uintptr {
	return uintptr((*value)(unsafe.Pointer(&v)).ptr)
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
		panic("no more pieces!")
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

func backupInstruction(location uintptr, minLen int) (code []byte, moreStackPtr uintptr) {
	moreStackPtr = uintptr(0)
	src := memoryAccess(location, 300)
	len := 0

	u := udisCodecInit(src)
	defer udisCodecFinal(u)

	var offset = 0
	var op *udisOP
	var adjustStack = uintptr(0)
	for offset < 200 {
		op = udisDecode(u)
		if op == nil {
			udisDisas(memoryAccess(location, 300))
			panic("backup instruction failed!")
		}
		if len < minLen {
			len += op.len
		}

		offset += op.len
		if op.ins == "jbe" {
			if op.jmpTo <= 0 {
				udisDisas(memoryAccess(location, 300))
				panic("jbe target not decoded!")
			}

			adjustStack = op.jmpTo
			break
		}

		if op.ins == "ret" || op.ins == "jmp" {
			break
		}
	}

	if adjustStack > 0 {
		udisCodecReset(u, memoryAccess(adjustStack, 50))
		op1 := udisDecode(u)
		if op1.ins != "call" {
			udisDisas(memoryAccess(location, 300))
			panic("wrong adjustStack addr!")
		}

		op2 := udisDecode(u)
		if op2.ins != "jmp" {
			udisDisas(memoryAccess(location, 300))
			panic("wrong adjustStack addr!")
		}

		if op2.jmpTo != location {
			udisDisas(memoryAccess(location, 300))
			panic("wrong adjustStack addr!")
		}

		moreStackPtr = adjustStack + uintptr(op1.len)
	}

	return src[0:len], moreStackPtr
}

func disasCode(bytes []byte) {
	udisDisas(bytes)
}
