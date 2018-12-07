package carrot

import "reflect"

type context struct {
	targetFuncBytes   []byte
	newFunc           *reflect.Value
	originalFuncBytes uintptr
}

// Patch is to patch function
//    targetFunc: func to replace
//    newFunc: new func
//    originalFunc: to recv original target func
func Patch(targetFunc, newFunc, originalFunc interface{}) bool {
	t := reflect.ValueOf(targetFunc)
	n := reflect.ValueOf(newFunc)
	o := reflect.ValueOf(originalFunc)

	checkType(t, n, o)

	return doPatch(t, n, o)
}

func checkType(t, n, o reflect.Value) {
	if t.Kind() != reflect.Func || n.Kind() != reflect.Func || o.Kind() != reflect.Func {
		panic("targetFunc, newFunc, originalFunc MUST be a func")
	}

	if t.Type() != n.Type() || t.Type() != o.Type() {
		panic("targetFunc, newFunc, originalFunc MUST be the same type")
	}
}

func doPatch(t, n, o reflect.Value) bool {
	return false
}
