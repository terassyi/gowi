package instruction

type I32Add struct{}

func (*I32Add) Opcode() Opcode {
	return I32_ADD
}

type I32Sub struct{}

func (*I32Sub) Opcode() Opcode {
	return I32_SUB
}

type I32Mul struct{}

func (*I32Mul) Opcode() Opcode {
	return I32_MUL
}
