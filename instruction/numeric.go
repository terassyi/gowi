package instruction

type I32Add struct{}

func (*I32Add) Opcode() Opcode {
	return I32_ADD
}

func (*I32Add) imm() any {
	return NoImm
}

type I32Sub struct{}

func (*I32Sub) Opcode() Opcode {
	return I32_SUB
}

func (*I32Sub) imm() any {
	return NoImm
}

type I32Mul struct{}

func (*I32Mul) Opcode() Opcode {
	return I32_MUL
}

func (*I32Mul) imm() any {
	return NoImm
}
