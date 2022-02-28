package instruction

type GetLocal struct {
	Imm uint32
}

func (*GetLocal) Opcode() Opcode {
	return GET_LOCAL
}

type SetLocal struct {
	Imm uint32
}

func (*SetLocal) Opcode() Opcode {
	return SET_LOCAL
}

type TeeLocal struct {
	Imm uint32
}

func (*TeeLocal) Opcode() Opcode {
	return TEE_LOCAL
}

type GetGlobal struct {
	Imm uint32
}

func (*GetGlobal) Opcode() Opcode {
	return GET_GLOBAL
}

type SetGlobal struct {
	Imm uint32
}

func (*SetGlobal) Opcode() Opcode {
	return SET_GLOBAL
}
