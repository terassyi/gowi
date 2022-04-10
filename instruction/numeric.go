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

type I32And struct{}

func (*I32And) Opcode() Opcode {
	return I32_AND
}

func (*I32And) imm() any {
	return NoImm
}

func (*I32And) String() string {
	return "i32.and"
}

type I64And struct{}

func (*I64And) Opcode() Opcode {
	return I64_AND
}

func (*I64And) imm() any {
	return NoImm
}

func (*I64And) String() string {
	return "i64.and"
}

type I32Or struct{}

func (*I32Or) Opcode() Opcode {
	return I32_OR
}

func (*I32Or) imm() any {
	return NoImm
}

func (*I32Or) String() string {
	return "i32.or"
}

type I64Or struct{}

func (*I64Or) Opcode() Opcode {
	return I64_OR
}

func (*I64Or) imm() any {
	return NoImm
}

func (*I64Or) String() string {
	return "i64.or"
}

type I32Xor struct{}

func (*I32Xor) Opcode() Opcode {
	return I32_XOR
}

func (*I32Xor) imm() any {
	return NoImm
}

func (*I32Xor) String() string {
	return "i32.xor"
}

type I64Xor struct{}

func (*I64Xor) Opcode() Opcode {
	return I64_XOR
}

func (*I64Xor) imm() any {
	return NoImm
}

func (*I64Xor) String() string {
	return "i64.xor"
}

type I32Shl struct{}

func (*I32Shl) Opcode() Opcode {
	return I32_SHL
}

func (*I32Shl) imm() any {
	return NoImm
}

func (*I32Shl) String() string {
	return "i32.shl"
}

type I64Shl struct{}

func (*I64Shl) Opcode() Opcode {
	return I64_SHL
}

func (*I64Shl) imm() any {
	return NoImm
}

func (*I64Shl) String() string {
	return "i64.shl"
}

type I32ShrU struct{}

func (*I32ShrU) Opcode() Opcode {
	return I32_SHR_U
}

func (*I32ShrU) imm() any {
	return NoImm
}

func (*I32ShrU) String() string {
	return "i32.shr_u"
}

type I64ShrU struct{}

func (*I64ShrU) Opcode() Opcode {
	return I64_SHR_U
}

func (*I64ShrU) imm() any {
	return NoImm
}

func (*I64ShrU) String() string {
	return "i64.shr_u"
}

type I32ShrS struct{}

func (*I32ShrS) Opcode() Opcode {
	return I32_SHR_S
}

func (*I32ShrS) imm() any {
	return NoImm
}

func (*I32ShrS) String() string {
	return "i32.shr_s"
}

type I64ShrS struct{}

func (*I64ShrS) Opcode() Opcode {
	return I64_SHR_S
}

func (*I64ShrS) imm() any {
	return NoImm
}

func (*I64ShrS) String() string {
	return "i64.shr_s"
}

type I32RotR struct{}

func (*I32RotR) Opcode() Opcode {
	return I32_ROTR
}

func (*I32RotR) imm() any {
	return NoImm
}

func (*I32RotR) String() string {
	return "i32.rotr"
}

type I64RotR struct{}

func (*I64RotR) Opcode() Opcode {
	return I64_ROTR
}

func (*I64RotR) imm() any {
	return NoImm
}

func (*I64RotR) String() string {
	return "i64.rotr"
}

type I32RotL struct{}

func (*I32RotL) Opcode() Opcode {
	return I32_ROTL
}

func (*I32RotL) imm() any {
	return NoImm
}

func (*I32RotL) String() string {
	return "i32.rotl"
}

type I64RotL struct{}

func (*I64RotL) Opcode() Opcode {
	return I64_ROTL
}

func (*I64RotL) imm() any {
	return NoImm
}

func (*I64RotL) String() string {
	return "i64.rotl"
}
