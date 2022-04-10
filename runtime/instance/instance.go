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
	~uint32 | ~uint64 | ~int32 | ~int64 | ~float32 | ~float64 | *Function | value.Vector
}

func GetVal[T ValueTypeSet](v value.Value) T {
	return v.(T)
}

type ExternalValueType uint8

const (
	ExternalValueTypeFunc   ExternalValueType = 0
	ExternalValueTypeTable  ExternalValueType = 1
	ExternalValueTypeMem    ExternalValueType = 2
	ExternalValueTypeGlobal ExternalValueType = 3
)

type ExternalValue interface {
	ExternalValueType() ExternalValueType
}

type ExternalValueTypeSet interface {
	*Function | *Table | *Memory | *Global
}

func GetExternVal[T ExternalValueTypeSet](v ExternalValue) T {
	return v.(T)
}
