package carrot

/*
#include <stdlib.h>
#include "thirdparty/udis86/libudis86/decode.c"
#include "thirdparty/udis86/libudis86/itab.c"
#include "thirdparty/udis86/libudis86/udis86.c"
#include "thirdparty/udis86/libudis86/syn.c"
#include "thirdparty/udis86/libudis86/syn-att.c"

ud_t* init() {
	return (ud_t*)malloc(sizeof(ud_t));
}

void Disas(char* code, int len)
{
	ud_t u;
	ud_init(&u);
	ud_set_mode(&u, 64);
	ud_set_input_buffer(&u, (uint8_t*)code, len);
	ud_set_syntax(&u, UD_SYN_ATT);
}
*/
import "C"
import "unsafe"

func disas(bytes []byte) {
	var code = (*C.uint8_t)(unsafe.Pointer(&bytes[0]))
	var length = (C.size_t)(len(bytes))
	var md = (C.uint8_t)(instructionLen())

	u := C.init()
	defer C.free(unsafe.Pointer(u))

	C.ud_init(u)
	C.ud_set_mode(u, md)
	C.ud_set_input_buffer(u, code, length)
}
