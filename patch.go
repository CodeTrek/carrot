package carrot

import (
	"reflect"
	"unsafe"
)

type patchContext struct {
	targetBytes   []byte
	originalBytes []byte
	bridge        []byte

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

func unpatch(t reflect.Value) {
	p, ok := patched[t.Pointer()]
	if !ok {
		return
	}

	doUnpatch(t.Pointer(), p)
}

func patch(t, r, o reflect.Value) bool {
	jmp2r := jmpTo(location(r))
	bridge := allocBridgePiece()
	bridgePtr := uintptr(unsafe.Pointer(&bridge[0]))
	backup, _ := backupInstruction(t.Pointer(), len(jmp2r))
	if len(bridge) < len(backup) {
		panic("bridge piece too small")
	}
	jmp2t := jmpTo(location(t) + uintptr(len(backup)))

	copyToLocation(bridgePtr, backup)
	copyToLocation(bridgePtr+uintptr(len(backup)), jmp2t)

	jmp2b := jmpTo(bridgePtr)
	originalBytes := make([]byte, len(jmp2b))
	copy(originalBytes, memoryAccess(location(o), len(jmp2b)))

	copyToLocation(t.Pointer(), jmp2r)
	copyToLocation(o.Pointer(), jmp2b)

	return true
}

func unpatchAll() {
	for t, p := range patched {
		doUnpatch(t, p)
	}
}

func doUnpatch(t uintptr, p patchContext) {
	delete(patched, t)
	delete(origins, (*p.replacement).Pointer())
}

func copyToLocation(location uintptr, data []byte) {
	g := makeWritable(location, len(data))
	defer g.restore()

	f := memoryAccess(location, len(data))
	copy(f, data[:])
}
