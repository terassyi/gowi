package instance

import (
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/value"
)

func evaluateConstInstr(instr instruction.Instruction) (value.Value, error) {
	switch instr.Opcode() {
	case instruction.I32_CONST:
		return value.I32(instruction.Imm[int32](instr)), nil
	case instruction.I64_CONST:
		return value.I64(instruction.Imm[int64](instr)), nil
	case instruction.F32_CONST:
		return value.F32(value.Float32FromUint32(instruction.Imm[uint32](instr))), nil
	case instruction.F64_CONST:
		return value.F64(value.Float64FromUint64(instruction.Imm[uint64](instr))), nil
	default:
		return nil, fmt.Errorf("evaluateConstInstr: %w: opcode=%x", instruction.NotConstInstruction, instr.Opcode())
	}
}

type ReferenceTypeSet interface {
	*Function
}

func GetRef[T ReferenceTypeSet](r value.Reference) T {
	return r.(T)
}

type ValueTypeSet interface {
	~int32 | ~int64 | ~float32 | ~float64 | *Function | value.Vector
}

func GetVal[T ValueTypeSet](v value.Value) T {
	return v.(T)
}

type ExternValueType uint8

const (
	ExternValTypeFunc   ExternValueType = 0
	ExternValTypeTable  ExternValueType = 1
	ExternValTypeMem    ExternValueType = 2
	ExternValTypeGlobal ExternValueType = 3
)

type ExternalValue interface {
	ExternalValueType() ExternValueType
}

type ExternValueTypeSet interface {
	*Function | *Table | *Memory | *Global
}

func GetExternVal[T ExternValueTypeSet](v ExternalValue) T {
	return v.(T)
}
