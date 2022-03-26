package instruction

type I32Add struct{}

func (*I32Add) Opcode() Opcode {
	return I32_ADD
}

func (*I32Add) imm() any {
	return NoImm
}

func (*I32Add) String() string {
	return "i32.add"
}

type I32Sub struct{}

func (*I32Sub) Opcode() Opcode {
	return I32_SUB
}

func (*I32Sub) imm() any {
	return NoImm
}

func (*I32Sub) String() string {
	return "i32.sub"
}

type I32Mul struct{}

func (*I32Mul) Opcode() Opcode {
	return I32_MUL
}

func (*I32Mul) imm() any {
	return NoImm
}

func (*I32Mul) String() string {
	return "i32.mul"
}
