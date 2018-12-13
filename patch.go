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
	backup, incrStack, targetCopiedLen, incrStkPtr, reachFuncEnd := backupInstruction(t.Pointer(), len(jmp2r))

	incrStkOffset := bridgeIncrStackOffset()
	if len(backup)+len(jmp2r) >= len(bridgePiece)-incrStkOffset {
		panic("bridge piece code section too small")
	}

	bridge := make([]byte, len(bridgePiece))
	copy(bridge[:incrStkOffset], backup)

	if !reachFuncEnd {
		jmp2t := jmpTo(getPtr(t) + uintptr(targetCopiedLen))
		copy(bridge[len(backup):incrStkOffset], jmp2t)
	}

	jmp2b := jmpTo(uintptr(unsafe.Pointer(&bridgePiece[0])))
	originalBytes := make([]byte, len(jmp2b))
	copy(originalBytes, memoryAccess(o.Pointer(), len(jmp2b)))

	targetBytes := make([]byte, len(jmp2r))
	copy(targetBytes, memoryAccess(t.Pointer(), len(jmp2r)))

	copyToLocation(t.Pointer(), jmp2r)
	copyToLocation(o.Pointer(), jmp2b)

	var targetAdjustBytes []byte
	if incrStkPtr > 0 {
		jmp2i := jmpTo(bridgePiecePtr + uintptr(incrStkOffset))
		targetAdjustBytes = make([]byte, len(incrStack)+len(jmp2i))
		copy(targetAdjustBytes, memoryAccess(incrStkPtr, len(targetAdjustBytes)))

		copy(bridge[incrStkOffset:], incrStack)
		copy(bridge[incrStkOffset+len(incrStack):], jmp2b)
		copyToLocation(incrStkPtr, jmp2i)
	}

	copyToLocation(bridgePiecePtr, bridge)

	//	fmt.Println("bridge")
	//	udisDisas(bridgePiece)

	patched[t.Pointer()] = patchContext{targetBytes, originalBytes, targetAdjustBytes, incrStkPtr, &bridgePiece, &r, &o}
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
	delete(origins, (*p.original).Pointer())
}

func copyToLocation(location uintptr, data []byte) {
	g := makeWritable(location, len(data))
	defer g.restore()

	f := memoryAccess(location, len(data))
	copy(f, data[:])
}

func getBridge(t reflect.Value) []byte {
	if p, ok := patched[t.Pointer()]; ok {
		return *(p.bridgeBytes)
	}

	return nil
}
