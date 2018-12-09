package carrot

import (
	"reflect"
	"sync"
)

var lock = sync.Mutex{}

// Patch is to patch function
//    targetFunc: func to replace
//    newFunc: new func
//    originalFunc: to recv original target func
func Patch(targetFunc, newFunc, originalFunc interface{}) bool {
	t := reflect.ValueOf(targetFunc)
	n := reflect.ValueOf(newFunc)
	o := reflect.ValueOf(originalFunc)

	checkType(t, n, o)

	lock.Lock()
	defer lock.Unlock()

	if isPatched(t) || isPatched(o) {
		return false
	}

	return patch(t, n, o)
}

// IsPatched to test wether f is patched
func IsPatched(f interface{}) bool {
	lock.Lock()
	defer lock.Unlock()

	return isPatched(reflect.ValueOf(f))
}

// Unpatch target
func Unpatch(target interface{}) {
	lock.Lock()
	defer lock.Unlock()

	t := reflect.ValueOf(target)
	if isPatched(t) {
		return
	}

	unpatch(t)
}

// UnpatchAll func
func UnpatchAll() {
	lock.Lock()
	defer lock.Unlock()

	unpatchAll()
}
