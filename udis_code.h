#ifndef _UDIS_CODE_H
#define _UDIS_CODE_H

#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "thirdparty/udis86/libudis86/types.h"
#include "thirdparty/udis86/libudis86/extern.h"
#include "thirdparty/udis86/libudis86/syn.h"
#include "thirdparty/udis86/libudis86/udint.h"

struct udis_backup_instr {
	int success;
	int reach_end;
	uintptr_t incr_stack_ptr;

	int copied_src_len;
	int copied_len;
	int incr_stack_len;
	uint8_t copied[128];
	uint8_t incr_stack[32];
};

struct udis_copy_instruction {
    const uint8_t* pc;
    uint8_t* buf;
};

typedef struct udis_backup_instr udis_backup_instr_t;
typedef struct udis_copy_instruction udis_copy_instruction_t;

#endif