#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "thirdparty/udis86/libudis86/types.h"
#include "thirdparty/udis86/libudis86/extern.h"
#include "thirdparty/udis86/libudis86/syn.h"

#ifndef BRIDGE_PIECE_SIZE
#	define BRIDGE_PIECE_SIZE 64
#endif

#ifndef BRIDGE_PAGE_COUNT
#	define BRIDGE_PAGE_COUNT 9
#endif

struct udis_codec {
	ud_t ud;
	const uint8_t* code;
	size_t len;
};

struct bridge {
	uintptr_t ptr;
	int       len;
};

struct udis_op {
	size_t      len;
	const char* ins;
	uintptr_t   jmp_to;
};

struct udis_code {
	uintptr_t ptr;
	int       len;
};

typedef struct bridge bridge_t;
typedef struct udis_codec  udis_codec_t;
typedef struct udis_op udis_op_t;
typedef struct udis_code udis_code_t;

static char bridge_buffer[4096 * BRIDGE_PAGE_COUNT] = {0};

static bridge_t udis_init_bridge()
{
	memset(bridge_buffer, 0xcc, sizeof(bridge_buffer));
	bridge_t b = { (uintptr_t)bridge_buffer, sizeof(bridge_buffer) };
	return b;
}

static size_t udis_bridge_piece_size()
{
	return BRIDGE_PIECE_SIZE;
}

static void udis_codec_reset(udis_codec_t* u, const uint8_t* code, size_t len)
{
	ud_init(&u->ud);
	ud_set_mode(&u->ud, sizeof(char*) * 8);
	ud_set_syntax(&u->ud, UD_SYN_INTEL);
	ud_set_input_buffer(&u->ud, code, len);
	ud_set_pc(&u->ud, (uint64_t)code);
	u->code = code;
	u->len = len;
}

static udis_codec_t* udis_codec_init(const uint8_t* code, size_t len)
{
	udis_codec_t* u = (udis_codec_t*)malloc(sizeof(udis_codec_t));
	udis_codec_reset(u, code, len);
	return u;
}

static void udis_codec_final(udis_codec_t* u)
{
	free((void*)u);
}

static udis_op_t udis_decode(udis_codec_t *u)
{
	udis_op_t d = {
		ud_decode(&u->ud),
		ud_lookup_mnemonic(ud_insn_mnemonic(&u->ud)),
		0
	};

	const struct ud_operand* opr = ud_insn_opr(&u->ud, 0);
	switch (ud_insn_mnemonic(&u->ud)) {
	case UD_Ijbe:
	case UD_Ijmp:
		if (opr->type == UD_OP_IMM || opr->type == UD_OP_JIMM) {
			d.jmp_to = ud_syn_rel_target(&u->ud, (struct ud_operand*)opr);
		}
		break;
	default:;
	}

	return d;
}

static udis_code_t udis_current_code(udis_codec_t* u)
{
	udis_code_t uc = { (uintptr_t)ud_insn_off(&u->ud), ud_insn_len(&u->ud) };
	return uc;
}

static void udis_print(udis_codec_t* u)
{
	int len = 0;
	while(len = ud_disassemble(&u->ud)) {
		if (ud_insn_mnemonic(&u->ud) == UD_Iint3) {
			break;
		}

		const char* ins = ud_insn_asm(&u->ud);
		char* addr = (char*)ud_insn_off(&u->ud);
        printf("0x%08x %-30s //", addr, ins);
		while (len-- > 0) {
			printf(" %02x", (uint8_t)(*addr++));
		}
		printf("\n");
    }
}
