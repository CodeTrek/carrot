package carrot

/*
#cgo CFLAGS: -DBRIDGE_PIECE_SIZE=32
#cgo CFLAGS: -DBRIDGE_PAGE_COUNT=5
#include "thirdparty/udis86/libudis86/decode.c"
#include "thirdparty/udis86/libudis86/itab.c"
#include "thirdparty/udis86/libudis86/udis86.c"
#include "thirdparty/udis86/libudis86/syn.c"
#include "thirdparty/udis86/libudis86/syn-intel.c"
//#include "thirdparty/udis86/libudis86/syn-att.c"
#include "udis_code.c"
*/
import "C"
import "unsafe"

type udisOP struct {
	len   int
	ins   string
	jmpTo uintptr
}

func udisInitBridge() (ptr uintptr, len uintptr) {
	bridge := C.udis_init_bridge()
	ptr = uintptr(bridge.ptr)
	len = uintptr(bridge.len)
	return
}

func udisCodecInit(data []byte) *C.udis_codec_t {
	var code = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	var length = (C.size_t)(len(data))

	return C.udis_codec_init(code, length)
}

func udisCodecReset(u *C.udis_codec_t, data []byte) {
	var code = (*C.uint8_t)(unsafe.Pointer(&data[0]))
	var length = (C.size_t)(len(data))

	C.udis_codec_reset(u, code, length)
}

func udisCodecFinal(u *C.udis_codec_t) {
	C.udis_codec_final(u)
}

func udisDisas(data []byte) {
	u := udisCodecInit(data)
	C.udis_print(u)
}

func udisBridgePieceSize() int {
	return int(C.udis_bridge_piece_size())
}

func udisDecode(u *C.udis_codec_t) *udisOP {
	c := C.udis_decode(u)
	var op = udisOP{int(c.len), C.GoString(c.ins), uintptr(c.jmp_to)}
	if op.ins == "invalid" {
		return nil
	}
	return &op
}

func udisCurrentCode(u *C.udis_codec_t) []byte {
	c := C.udis_current_code(u)
	ptr := uintptr(c.ptr)
	len := int(c.len)

	if len == 0 {
		return nil
	}

	return memoryAccess(ptr, len)
}
