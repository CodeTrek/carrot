package carrot

import (
	"reflect"
)

type patchContext struct {
	targetBytes   []byte
	originalBytes []byte
	piceBytes     []byte

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
	disas(memoryAccess(o.Pointer(), 50))

	jmp2r := jmpTo(location(r))
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
