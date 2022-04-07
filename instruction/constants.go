package instruction

import "github.com/terassyi/gowi/types"

type I32Const struct {
	Imm int32
}

func (*I32Const) Opcode() Opcode {
	return I32_CONST
}

func (i32c *I32Const) imm() any {
	return i32c.Imm
}

func (*I32Const) String() string {
	return "i32.const"
}

type I64Const struct {
	Imm int64
}

func (*I64Const) Opcode() Opcode {
	return I64_CONST
}

func (i64c *I64Const) imm() any {
	return i64c.Imm
}

func (*I64Const) String() string {
	return "i64.const"
}

type F32Const struct {
	Imm uint32
}

func (*F32Const) Opcode() Opcode {
	return F32_CONST
}

func (f32c *F32Const) imm() any {
	return f32c.Imm
}

func (*F32Const) String() string {
	return "f32.const"
}

type F64Const struct {
	Imm uint64
}

func (*F64Const) Opcode() Opcode {
	return F64_CONST
}

func (f64c *F64Const) imm() any {
	return f64c.Imm
}

func (*F64Const) String() string {
	return "f64.const"
}

func IsConst(instr Instruction) bool {
	switch instr.Opcode() {
	case I32_CONST:
		return true
	case I64_CONST:
		return true
	case F32_CONST:
		return true
	case F64_CONST:
		return true
	default:
		return false
	}
}

func GetConstType(instr Instruction) (types.ValueType, error) {
	switch instr.Opcode() {
	case I32_CONST:
		return types.I32, nil
	case I64_CONST:
		return types.I64, nil
	case F32_CONST:
		return types.F32, nil
	case F64_CONST:
		return types.F64, nil
	default:
		return types.ValueType(0xff), NotConstInstruction
	}
}
