package carrot

import (
	"reflect"
	"sync"
)

type patchContext struct {
	isTarget          bool // true - targetFunc    false - originalFunc
	targetFuncBytes   []byte
	originalFuncBytes []byte
	newFunc           *reflect.Value
	originalFunc      *reflect.Value
}

var (
	lock = sync.Mutex{}

	patched = make(map[uintptr]patchContext)
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
	_, ok := patched[t.Pointer()]
	return ok
}

func doPatch(t, n, o reflect.Value) bool {
	lock.Lock()
	defer lock.Unlock()

	disas(memoryAccess(t.Pointer(), 200))
	return false
}

func unPatch(t reflect.Value) {
	lock.Lock()
	defer lock.Unlock()

	p, ok := patched[t.Pointer()]
	if !ok {
		return
	}

	delete(patched, t.Pointer())
	delete(patched, (*p.originalFunc).Pointer())
}
