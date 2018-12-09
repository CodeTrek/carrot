package carrot

import (
	"reflect"
)

type patchContext struct {
	targetBytes []byte
	bridgeBytes []byte
	piceBytes   []byte

	replacement *reflect.Value
	bridge      *reflect.Value
}

var (
	patched = make(map[uintptr]patchContext)
	origins = make(map[uintptr]bool)
)

func checkType(t, r, b reflect.Value) {
	if t.Kind() != reflect.Func || r.Kind() != reflect.Func || b.Kind() != reflect.Func {
		panic("target, replacement, bridge MUST be a func")
	}

	if t.Type() != r.Type() || t.Type() != b.Type() {
		panic("target, replacement, bridge MUST be the same type")
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

func patch(t, r, b reflect.Value) bool {
	//	disas(memoryAccess(b.Pointer(), 50))

	jmp2r := jmpTo(getFuncAddr(r))
	//	pices := allocPices()

	copyToLocation(t.Pointer(), jmp2r)
	return false
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
