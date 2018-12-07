package carrot

/*
#include "thirdparty/udis86/libudis86/decode.c"
#include "thirdparty/udis86/libudis86/itab.c"
#include "thirdparty/udis86/libudis86/udis86.c"
#include "thirdparty/udis86/libudis86/syn.c"
#include "thirdparty/udis86/libudis86/syn-att.c"

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
