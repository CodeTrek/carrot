package carrot

import (
	"reflect"
)

type patchContext struct {
	targetFuncBytes   []byte
	originalFuncBytes []byte
	newFunc           *reflect.Value
	originalFunc      *reflect.Value
}

var (
	patched       = make(map[uintptr]patchContext)
	origins       = make(map[uintptr]bool)
	freeCodePices = initCodePage(4)
	usedCodePices = make(map[uintptr][]byte)
)

func checkType(t, n, o reflect.Value) {
	if t.Kind() != reflect.Func || n.Kind() != reflect.Func || o.Kind() != reflect.Func {
		panic("targetFunc, newFunc, originalFunc MUST be a func")
	}

	if t.Type() != n.Type() || t.Type() != o.Type() {
		panic("targetFunc, newFunc, originalFunc MUST be the same type")
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

func patch(t, n, o reflect.Value) bool {
	//	disas(memoryAccess(o.Pointer(), 50))
	return false
}

func unpatchAll() {
	for t, p := range patched {
		doUnpatch(t, p)
	}
}

func doUnpatch(t uintptr, p patchContext) {
	delete(patched, t)
	delete(origins, (*p.originalFunc).Pointer())
}
