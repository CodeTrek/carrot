#include "udis_code.h"

/*
x64
0x495b20 lea rax, [rip+0x5a3f8]         // 48 8d 05 f8 a3 05 00
0x495b27 lea rbx, [rip+0x5a3f8]         // 48 8d 1d f8 a3 05 00
0x495b2e lea rcx, [rip+0x5a3f8]         // 48 8d 0d f8 a3 05 00
0x495b35 lea rdx, [rip+0x5a3f8]         // 48 8d 15 f8 a3 05 00
i386
0x495b3c lea eax, [rip+0x5a3f8]         //    8d 05 f8 a3 05 00
0x495b42 lea ebx, [rip+0x5a3f8]         //    8d 1d f8 a3 05 00
0x495b48 lea rcx, [rip+0x5a3f8]         //    8d 0d f8 a3 05 00
0x495b4e lea r10, [rip+0x5a3f8]         //    8d 15 f8 a3 05 00

x64
0x495b54 mov rax, 0x7ff48e490           // 48 b8 90 e4 48 ff 07 00 00 00
0x495b5e mov rbx, 0x7ff48e490           // 48 bb 90 e4 48 ff 07 00 00 00
0x495b68 mov rcx, 0x7ff48e490           // 48 b9 90 e4 48 ff 07 00 00 00
0x495b72 mov rdx, 0x7ff48e490           // 48 ba 90 e4 48 ff 07 00 00 00
i386
0x495b7c mov eax, 0xff48e490            //    b8 90 e4 48 ff
0x495b81 mov ebx, 0xff48e490            //    bb 90 e4 48 ff
0x495b86 mov ecx, 0xff48e490            //    b9 90 e4 48 ff
0x495b8b mov edx, 0xff48e490            //    ba 90 e4 48 ff
*/

static int copy_raw(ud_t* u, udis_copy_instruction_t* p)
{
	int len = ud_insn_len(u);
	memcpy(p->buf, p->pc, len);
	return len;
}

static int copy_lea(ud_t* u, udis_copy_instruction_t* p)
{
    // lea rax|rbx|rcx|rdx, [rip+0x5a588] => mov rax|rbx|rcx|rbx, 0x5a588

    int offset = 0;
    if (
#ifndef _M_IX86
        p->pc[offset++] == 0x48 &&
#endif
        p->pc[offset++] == 0x8d
    ) {
        int d = 0;
        switch (p->pc[offset++]) {
        case 0x05: d = 0xb8; break;
        case 0x1d: d = 0xbb; break;
        case 0x0d: d = 0xb9; break;
        case 0x15: d = 0xba; break;
        }
        if (d > 0) {
            uintptr_t addr = (uintptr_t)p->pc + ud_insn_len(u) + *(uint32_t*)(p->pc + offset);
            uint8_t ins[] = {
#ifndef _M_IX86
                0x48,
#endif
                0xb8,
                (uint8_t)(addr & 0xff),
                (uint8_t)((addr >> 8) & 0xff),
                (uint8_t)((addr >> 16) & 0xff),
                (uint8_t)((addr >> 24) & 0xff),
#ifndef _M_IX86
                (uint8_t)((addr >> 32) & 0xff),
                (uint8_t)((addr >> 40) & 0xff),
                (uint8_t)((addr >> 48) & 0xff),
                (uint8_t)((addr >> 56) & 0xff),
#endif
            };

            memcpy(p->buf, ins, sizeof(ins));
            return sizeof(ins);
        }
    }

    return copy_raw(u, p);
}

static int copy_instruction(ud_t* u, udis_copy_instruction_t* p)
{
    switch (ud_insn_mnemonic(u)) {
    case UD_Ilea:
        return copy_lea(u, p);
    }

    return copy_raw(u, p);
}
