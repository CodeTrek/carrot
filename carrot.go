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

// Patch is to patch function
//    targetFunc: func to replace
//    newFunc: new func
//    originalFunc: to recv original target func
func Patch(targetFunc, newFunc, originalFunc interface{}) bool {
	t := reflect.ValueOf(targetFunc)
	n := reflect.ValueOf(newFunc)
	o := reflect.ValueOf(originalFunc)

	checkType(t, n, o)

	if isPatched(t) || isPatched(o) {
		return false
	}

	return doPatch(t, n, o)
}

// IsPatched to test wether f is patched
func IsPatched(f interface{}) bool {
	return isPatched(reflect.ValueOf(f))
}

// Unpatch target
func Unpatch(target interface{}) {
	t := reflect.ValueOf(target)
	if isPatched(t) {
		return
	}

}

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
	disas(memoryAccess(t.Pointer(), 200))
	return false
}

func unPatch(t reflect.Value) {
	p, ok := patched[t.Pointer()]
	if !ok {
		return
	}

	delete(patched, t.Pointer())
	delete(patched, (*p.originalFunc).Pointer())
}
