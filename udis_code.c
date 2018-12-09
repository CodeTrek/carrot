#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "thirdparty/udis86/libudis86/types.h"
#include "thirdparty/udis86/libudis86/extern.h"

struct udis {
	ud_t ud;
	uint8_t* code;
	size_t len;
};

typedef struct udis  udis_t;

static udis_t* udis_init(const uint8_t* code, size_t len) {
	udis_t* u = (udis_t*)malloc(sizeof(udis_t));
	ud_init(&u->ud);
	ud_set_mode(&u->ud, sizeof(char*) * 8);
	ud_set_syntax(&u->ud, UD_SYN_INTEL);
	ud_set_input_buffer(&u->ud, code, len);
	ud_set_pc(&u->ud, (uint64_t)code);

	return u;
}

static void udis_final(udis_t* u) {
	free((void*)u);
}

static void udis_print(udis_t* u) {
	int len = 0;
	while(len = ud_disassemble(&u->ud)) {
		if (ud_insn_mnemonic(&u->ud) == UD_Iint3) {
			break;
		}

		const char* ins = ud_insn_asm(&u->ud);
		char* addr = (char*)ud_insn_off(&u->ud);
        printf("0x%016X:%-2d %-32s", addr, len, ins);
		while (len-- > 0) {
			printf("%02X ", (uint8_t)(*addr++));
		}
		printf("\n");
    }
}
