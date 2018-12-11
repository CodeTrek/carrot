// +build windows

package carrot

import (
	"syscall"
	"unsafe"
)

var procVirtualProtect = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualProtect")

func vProtect(location uintptr, len int, attrib uint32) uint32 {
	var tmp uint32
	ret, _, _ := procVirtualProtect.Call(
		location,
		uintptr(len),
		uintptr(attrib),
		uintptr(unsafe.Pointer(&tmp)))
	if ret == 0 {
		panic(syscall.GetLastError())
	}

	return tmp
}

var vAlloc = syscall.NewLazyDLL("kernel32.dll").NewProc("VirtualAlloc")

func allocPage() []byte {
	// MEM_COMMIT = 0x00001000
	// PAGE_EXECUTE_READ = 0x20
	ret, _, _ := vAlloc.Call(
		uintptr(0),
		uintptr(syscall.Getpagesize()),
		uintptr(0x00001000),
		uintptr(0x20))
	if ret == 0 {
		panic(syscall.GetLastError())
	}

	return memoryAccess(uintptr(ret), syscall.Getpagesize())
}

type memProtectGuard struct {
	oldPerms uint32
	location uintptr
	len      int
}

func makeWritable(location uintptr, len int) *memProtectGuard {
	// PAGE_EXECUTE_READWRITE = 0x40
	var oldPerms = vProtect(location, len, 0x40)
	return &memProtectGuard{oldPerms, location, len}
}

func (g *memProtectGuard) restore() {
	vProtect(g.location, g.len, g.oldPerms)
}
