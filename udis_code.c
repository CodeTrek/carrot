#include "udis_code.h"

extern int copy_instruction(ud_t* u, udis_copy_instruction_t* p);

static void udis_init(ud_t* u, const uint8_t* code, size_t len)
{
	ud_init(u);
	ud_set_mode(u, sizeof(char*) * 8);
	ud_set_syntax(u, UD_SYN_ATT);
	ud_set_input_buffer(u, code, len);
	ud_set_pc(u, (uint64_t)code);
}

static void udis_disas(const uint8_t* code, size_t code_len)
{
	ud_t u;
	udis_init(&u, code, code_len);

	int len = 0;
	while(len = ud_disassemble(&u)) {
		if (ud_insn_mnemonic(&u) == UD_Iint3) {
			break;
		}

		const char* ins = ud_insn_asm(&u);
		char* addr = (char*)ud_insn_off(&u);
        printf("0x%06" FMT64 "x %-30s //", addr, ins);
		while (len-- > 0) {
			printf(" %02x", (uint8_t)(*addr++));
		}
		printf("\n");
    }
}

static int64_t udis_real_target(ud_t* u, int n)
{
	const struct ud_operand* opr = ud_insn_opr(u, n);
	switch (ud_insn_mnemonic(u)) {
	case UD_Ijbe:
	case UD_Ijmp:
		if (opr->type == UD_OP_IMM || opr->type == UD_OP_JIMM) {
			return ud_syn_rel_target(u, (struct ud_operand*)opr);
		}
		break;
	default:;
		return -1;
	}
}

static udis_backup_instr_t udis_backup_instruction(const uint8_t* code, size_t len, size_t jmp_len, uintptr_t data_ptr)
{
	udis_backup_instr_t result;
	memset(&result, 0, sizeof(result));
	result.success = 1;

	ud_t u;
	udis_init(&u, code, len);

	int copied_src_len = 0;
	int offset = 0;
	while (offset < 200 && offset < len - 20) {
		int current_len = ud_decode(&u);
		enum ud_mnemonic_code current_ins = ud_insn_mnemonic(&u);
		if (current_ins == UD_Iinvalid) {
			result.success = 0;
			break;
		}

		if (copied_src_len < jmp_len) {
			udis_copy_instruction_t cp_ins = {
				data_ptr,
				code + copied_src_len,
				&result.copied[result.copied_len],
				result.data,
				&result.data_len
			};
			result.copied_len += copy_instruction(&u, &cp_ins);

			copied_src_len += current_len;

			if (current_ins == UD_Iret || current_ins == UD_Ijmp) {
				result.reach_end = 1;
				break;
			}
		}

		offset += current_len;
		if (current_ins == UD_Ijbe) {
			const struct ud_operand* opr = ud_insn_opr(&u, 0);
			if (opr->type == UD_OP_IMM || opr->type == UD_OP_JIMM) {
				result.adjust_stack_jmp = ud_syn_rel_target(&u, (struct ud_operand*)opr);
			} else {
				result.success = 0;
			}
			break;
		}

		if (current_ins == UD_Iret || current_ins == UD_Ijmp) {
			break;
		}
	}

	if (result.adjust_stack_jmp > 0) {
		udis_init(&u, (const uint8_t*)result.adjust_stack_jmp, 50);
		ud_decode(&u);

		if (ud_insn_mnemonic(&u) == UD_Icall) {
			int len = ud_insn_len(&u);
			ud_decode(&u);
			if (ud_insn_mnemonic(&u) != UD_Ijmp) {
				result.success = 0;
			} else {
				result.adjust_stack_jmp += len;
			}
		}
	}

	return result;
}
