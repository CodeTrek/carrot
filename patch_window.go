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

func makePageExecutable(location uintptr) {
	// PAGE_EXECUTE_READ = 0x20
	vProtect(location, syscall.Getpagesize(), 0x20)
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
