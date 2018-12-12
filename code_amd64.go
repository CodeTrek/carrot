package carrot

// Assembles a jump to a function value
func jmpTo(to uintptr) []byte {
	return []byte{
		0x48, 0xBA,
		byte(to),
		byte(to >> 8),
		byte(to >> 16),
		byte(to >> 24),
		byte(to >> 32),
		byte(to >> 40),
		byte(to >> 48),
		byte(to >> 56), // movabs rdx,to
		0xFF, 0xe2,     // jmp QWORD PTR [rdx]
	}
}

func movRAX(v uintptr) []byte {
	return []byte{
		0x48, 0xB8,
		byte(v),
		byte(v >> 8),
		byte(v >> 16),
		byte(v >> 24),
		byte(v >> 32),
		byte(v >> 40),
		byte(v >> 48),
		byte(v >> 56), // movabs rax,v
	}
}
