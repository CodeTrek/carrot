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
	uintptr_t adjust_stack_jmp;

	int copied_src_len;
	int data_len;
	int copied_len;
	uint8_t copied[128];
	uint8_t data[128];
};

struct udis_copy_instruction {
    uintptr_t raw_data_ptr;
    const uint8_t* pc;
    uint8_t* buf;
    uint8_t* data;
    int* data_used;
};

typedef struct udis_backup_instr udis_backup_instr_t;
typedef struct udis_copy_instruction udis_copy_instruction_t;

#endif