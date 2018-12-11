//+build !windows

package carrot

import (
	"syscall"
	"unsafe"
)

var pageSize = syscall.Getpagesize()

func mProtect(b []byte, prot int) {
	var _p0 = unsafe.Pointer(&b[0])
	_, _, e1 := syscall.RawSyscall(syscall.SYS_MPROTECT, uintptr(_p0), uintptr(len(b)), uintptr(prot))
	if e1 != 0 {
		panic(e1)
	}
	return
}

func allocPage() []byte {
	//f, err := os.OpenFile("/dev/zero", os.O_RDWR|os.O_CREATE, 0777)
	//	defer f.Close()

	data, err := syscall.Mmap(-1, 0, syscall.Getpagesize(),
		syscall.PROT_READ|syscall.PROT_EXEC, syscall.MAP_PRIVATE|syscall.MAP_ANONYMOUS)
	if err != nil {
		panic(err)
	}

	return data
}

type memProtectGuard struct {
	page     uintptr
	nextPage uintptr
}

func makeWritable(location uintptr, len int) *memProtectGuard {
	attrib := syscall.PROT_READ | syscall.PROT_WRITE | syscall.PROT_EXEC
	page := pageStart(location)
	nextPage := uintptr(0)

	mProtect(memoryAccess(page, pageSize), attrib)
	page2 := page + uintptr(pageSize)
	if location+uintptr(len) >= page2 {
		nextPage = page2
		mProtect(memoryAccess(page2, pageSize), attrib)
	}

	return &memProtectGuard{page, nextPage}
}

func (g *memProtectGuard) restore() {
	attrib := syscall.PROT_READ | syscall.PROT_EXEC
	mProtect(memoryAccess(g.page, pageSize), attrib)
	if g.nextPage > 0 {
		mProtect(memoryAccess(g.nextPage, pageSize), attrib)
	}
}
