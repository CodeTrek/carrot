package carrot

import (
	"reflect"
	"sync"
)

/*
void int3()
{
	__asm__("int $0x03");
}
*/
import "C"

// Break int3
func Break() {
	C.int3()
}

var lock = sync.Mutex{}

// Patch is to patch function
//    target: func to replace
//    replacement: new func
//    original: to recv original target func pointer
func Patch(target, replacement, original interface{}) bool {
	t := reflect.ValueOf(target)
	r := reflect.ValueOf(replacement)
	o := reflect.ValueOf(original)

	checkType(t, r, o)

	lock.Lock()
	defer lock.Unlock()

	if isPatched(t) || isPatched(o) {
		return false
	}

	return patch(t, r, o)
}

// IsPatched to test wether f is patched or used as bridge
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

// Disas function
func Disas(target interface{}) {
	t := reflect.ValueOf(target)
	if t.Kind() != reflect.Func {
		panic("f MUST BE func")
	}

	udisDisas(memoryAccess(t.Pointer(), 6000))
}
