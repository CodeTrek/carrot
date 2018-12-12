package carrot

import (
	"reflect"
	"unsafe"
)

type patchContext struct {
	targetBytes               []byte
	originalBytes             []byte
	targetAdjustStackRetBytes []byte
	targetAdjustStackRet      uintptr

	bridgeBytes *[]byte
	replacement *reflect.Value
	original    *reflect.Value
}

var (
	patched = make(map[uintptr]patchContext)
	origins = make(map[uintptr]bool)
)

func checkType(t, r, o reflect.Value) {
	if t.Kind() != reflect.Func || r.Kind() != reflect.Func || o.Kind() != reflect.Func {
		panic("target, replacement, original MUST be a func")
	}

	if t.Type() != r.Type() || t.Type() != o.Type() {
		panic("target, replacement, original MUST be the same type")
	}
}

func isPatched(t reflect.Value) bool {
	if _, ok := patched[t.Pointer()]; ok {
		return true
	}

	if _, ok := origins[t.Pointer()]; ok {
		return true
	}

	return false
}

func patch(t, r, o reflect.Value) bool {
	jmp2r := jmpTo(getPtr(r))
	bridgePiece := allocBridgePiece()
	bridgePiecePtr := uintptr(unsafe.Pointer(&bridgePiece[0]))
	backup, dataUsed, targetCopiedLen, moreStackJmp, reachFuncEnd :=
		backupInstruction(t.Pointer(), len(jmp2r), bridgePiecePtr+uintptr(bridgePieceDataOffset()))
	if len(dataUsed) > len(bridgePiece)-bridgePieceDataOffset() {
		panic("bridge piece data section too small")
	}
	dataOffset := bridgePieceDataOffset()
	if len(backup)+len(jmp2r) >= len(bridgePiece)-dataOffset {
		panic("bridge piece code section too small")
	}

	bridge := make([]byte, len(bridgePiece))
	copy(bridge[:dataOffset], backup)
	copy(bridge[dataOffset:], dataUsed)

	if !reachFuncEnd {
		jmp2t := jmpTo(getPtr(t) + uintptr(targetCopiedLen))
		copy(bridge[len(backup):dataOffset], jmp2t)
	}

	copyToLocation(bridgePiecePtr, bridge)

	jmp2b := jmpTo(uintptr(unsafe.Pointer(&bridgePiece[0])))
	originalBytes := make([]byte, len(jmp2b))
	copy(originalBytes, memoryAccess(o.Pointer(), len(jmp2b)))

	targetBytes := make([]byte, len(jmp2r))
	copy(targetBytes, memoryAccess(t.Pointer(), len(jmp2r)))

	copyToLocation(t.Pointer(), jmp2r)
	copyToLocation(o.Pointer(), jmp2b)

	var targetAdjustBytes []byte
	if moreStackJmp > 0 {
		targetAdjustBytes = make([]byte, len(jmp2b))
		copy(targetAdjustBytes, memoryAccess(moreStackJmp, len(jmp2b)))
		copyToLocation(moreStackJmp, jmp2b)
	}

	//	fmt.Println("bridge")
	//	udisDisas(bridgePiece)

	patched[t.Pointer()] = patchContext{targetBytes, originalBytes, targetAdjustBytes, moreStackJmp, &bridgePiece, &r, &o}
	origins[o.Pointer()] = true
	return true
}

func unpatch(t reflect.Value) {
	p, ok := patched[t.Pointer()]
	if !ok {
		return
	}

	doUnpatch(t.Pointer(), p)
}

func unpatchAll() {
	for t, p := range patched {
		doUnpatch(t, p)
	}
}

func doUnpatch(t uintptr, p patchContext) {
	copyToLocation(t, p.targetBytes)
	copyToLocation(p.original.Pointer(), p.originalBytes)
	if p.targetAdjustStackRet > 0 {
		copyToLocation(p.targetAdjustStackRet, p.targetAdjustStackRetBytes)
	}
	freeBridgePiece(*p.bridgeBytes)

	delete(patched, t)
	delete(origins, (*p.replacement).Pointer())
}

func copyToLocation(location uintptr, data []byte) {
	g := makeWritable(location, len(data))
	defer g.restore()

	f := memoryAccess(location, len(data))
	copy(f, data[:])
}
