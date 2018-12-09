//+build !windows

package carrot

import "syscall"

// this function is super unsafe
// aww yeah
// It copies a slice to a raw memory location, disabling all memory protection before doing so.
func copyToLocation(location uintptr, data []byte) {
	attrib := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	start := pageStart(location)
	page_size := syscall.Getpagesize()

	page := memoryAccess(start, page_size)
	err1 := syscall.Mprotect(page, attrib)
	if err1 != nil {
		panic(err1)
	}

	if start+uintptr(len(data)) > start+uintptr(page_size) {
		nextPage := memoryAccess(start+uintptr(page_size), page_size)
		err2 := syscall.Mprotect(nextPage, attrib)
		if err2 != nil {
			panic(err2)
		}
	}

	f := memoryAccess(location, len(data))
	copy(f, data[:])

	attrib = syscall.PROT_READ | syscall.PROT_EXEC
	err1 = syscall.Mprotect(page, attrib)
	if err1 != nil {
		panic(err1)
	}

	if start+uintptr(len(data)) > start+uintptr(page_size) {
		err2 = syscall.Mprotect(nextPage, attrib)
		if err2 != nil {
			panic(err2)
		}
	}
}
