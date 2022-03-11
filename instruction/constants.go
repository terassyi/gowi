package instruction

type I32Const struct {
	Imm int32
}

func (*I32Const) Opcode() Opcode {
	return I32_CONST
}

type I64Const struct {
	Imm int64
}

func (*I64Const) Opcode() Opcode {
	return I64_CONST
}

type F32Const struct {
	Imm uint32
}

func (*F32Const) Opcode() Opcode {
	return F32_CONST
}

type F64Const struct {
	Imm uint64
}

func (*F64Const) Opcode() Opcode {
	return F64_CONST
}
