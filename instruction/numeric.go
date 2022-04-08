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

type I32DivS struct{}

func (*I32DivS) Opcode() Opcode {
	return I32_DIV_S
}

func (*I32DivS) imm() any {
	return NoImm
}

func (*I32DivS) String() string {
	return "i32.div_s"
}

type I32DivU struct{}

func (*I32DivU) Opcode() Opcode {
	return I32_DIV_U
}

func (*I32DivU) imm() any {
	return NoImm
}

func (*I32DivU) String() string {
	return "i32.div_u"
}

type I64Add struct{}

func (*I64Add) Opcode() Opcode {
	return I64_ADD
}

func (*I64Add) imm() any {
	return NoImm
}

func (*I64Add) String() string {
	return "i64.add"
}

type I64Sub struct{}

func (*I64Sub) Opcode() Opcode {
	return I64_SUB
}

func (*I64Sub) imm() any {
	return NoImm
}

func (*I64Sub) String() string {
	return "i64.sub"
}

type I64Mul struct{}

func (*I64Mul) Opcode() Opcode {
	return I64_MUL
}

func (*I64Mul) imm() any {
	return NoImm
}

func (*I64Mul) String() string {
	return "i64.mul"
}

type I64DivS struct{}

func (*I64DivS) Opcode() Opcode {
	return I64_DIV_S
}

func (*I64DivS) imm() any {
	return NoImm
}

func (*I64DivS) String() string {
	return "i64.div_s"
}

type I64DivU struct{}

func (*I64DivU) Opcode() Opcode {
	return I64_DIV_U
}

func (*I64DivU) imm() any {
	return NoImm
}

func (*I64DivU) String() string {
	return "i64.div_u"
}

type I32RemS struct{}

func (*I32RemS) Opcode() Opcode {
	return I32_REM_S
}

func (*I32RemS) imm() any {
	return NoImm
}

func (*I32RemS) String() string {
	return "i32.rem_s"
}

type I32RemU struct{}

func (*I32RemU) Opcode() Opcode {
	return I32_REM_U
}

func (*I32RemU) imm() any {
	return NoImm
}

func (*I32RemU) String() string {
	return "i32.rem_u"
}

type I64RemS struct{}

func (*I64RemS) Opcode() Opcode {
	return I64_REM_S
}

func (*I64RemS) imm() any {
	return NoImm
}

func (*I64RemS) String() string {
	return "i64.rem_s"
}

type I64RemU struct{}

func (*I64RemU) Opcode() Opcode {
	return I64_REM_U
}

func (*I64RemU) imm() any {
	return NoImm
}

func (*I64RemU) String() string {
	return "i64.rem_u"
}
