//+build !windows

package carrot

import (
	"syscall"
)

func mProtect(location uintptr, attrib int) {
	err := syscall.Mprotect(memoryAccess(location, syscall.Getpagesize()), attrib)
	if err != nil {
		panic(err)
	}
}

func makePageExecutable(location uintptr) {
	mProtect(location, syscall.PROT_READ|syscall.PROT_EXEC)
}

type memProtectGuard struct {
	page     uintptr
	nextPage uintptr
}

func makeWritable(location uintptr, len int) *memProtectGuard {
	attrib := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	page := pageStart(location)
	nextPage := uintptr(0)

	mProtect(page, attrib)
	page2 := page + uintptr(syscall.Getpagesize())
	if location+uintptr(len) >= page2 {
		nextPage = page2
		mProtect(page2, attrib)
	}

	return &memProtectGuard{page, nextPage}
}

func (g *memProtectGuard) restore() {
	attrib := syscall.PROT_READ | syscall.PROT_EXEC
	mProtect(g.page, attrib)
	if g.nextPage > 0 {
		mProtect(g.nextPage, attrib)
	}
}
