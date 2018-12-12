package main

/*
#include "../thirdparty/udis86/libudis86/decode.c"
#include "../thirdparty/udis86/libudis86/itab.c"
#include "../thirdparty/udis86/libudis86/udis86.c"
#include "../thirdparty/udis86/libudis86/syn.c"
#include "../thirdparty/udis86/libudis86/syn-att.c"
#include "../udis_copy_instruction.c"
#include "../udis_code.c"
void asmdd()
{
	__asm__("lea 0x5a3f8(%rip), %rax");
	__asm__("lea 0x5a3f8(%rip), %rbx");
	__asm__("lea 0x5a3f8(%rip), %rcx");
	__asm__("lea 0x5a3f8(%rip), %rdx");

	__asm__("mov $0x7ff48e490, %rax");
	__asm__("mov $0x7ff48e490, %rbx");
	__asm__("mov $0x7ff48e490, %rcx");
	__asm__("mov $0x7ff48e490, %rdx");

	__asm__("lea 0x5a3f8(%rip), %eax");
	__asm__("lea 0x5a3f8(%rip), %ebx");
	__asm__("lea 0x5a3f8(%rip), %ecx");
	__asm__("lea 0x5a3f8(%rip), %edx");

	__asm__("mov $0xff48e490, %eax");
	__asm__("mov $0xff48e490, %ebx");
	__asm__("mov $0xff48e490, %ecx");
	__asm__("mov $0xff48e490, %edx");
}
static void dis() {
	udis_disas((const uint8_t*)asmdd, 200);
}
*/
import "C"

func main() {
	C.dis()
}
