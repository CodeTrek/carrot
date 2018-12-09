package carrot

import (
	"reflect"
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
