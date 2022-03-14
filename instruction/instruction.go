package instruction

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/terassyi/gowi/types"
)

var (
	InvalidOpcode       error = errors.New("Invalid opcode")
	NotImplemented      error = errors.New("Not implemented")
	NotConstInstruction error = errors.New("Not const instruction")
)

type Instruction interface {
	Opcode() Opcode
}

func New(opcode uint8) (Instruction, error) {
	return nil, nil
}

func Decode(buf *bytes.Buffer) (Instruction, error) {
	opcode, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("Instruction decode opcode: %w", err)
	}
	switch Opcode(opcode) {
	case UNREACHABLE:
		return &Unreachable{}, nil
	case NOP:
		return &Nop{}, nil
	case BLOCK:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(block) decode: %w", err)
		}
		return &Block{Imm: types.BlockType(imm)}, nil
	case LOOP:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(loop) decode: %w", err)
		}
		return &Loop{Imm: types.BlockType(imm)}, nil
	case IF:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(if) decode: %w", err)
		}
		return &If{Imm: types.BlockType(imm)}, nil
	case ELSE:
		return &Else{}, nil
	case END:
		return &End{}, nil
	case BR:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(br) decode: %w", err)
		}
		return &Br{Imm: uint32(imm)}, nil
	case BR_IF:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(br_if) decode: %w", err)
		}
		return &BrIf{Imm: uint32(imm)}, nil
	case BR_TABLE:
		count, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(br_table) decode count: %w", err)
		}
		targets := make([]uint32, 0, int(count))
		for i := 0; i < int(count); i++ {
			t, _, err := types.DecodeVarUint32(buf)
			if err != nil {
				return nil, fmt.Errorf("Instruction(br_table) decode target: %w", err)
			}
			targets = append(targets, uint32(t))
		}
		def, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(br_table) decode default: %w", err)
		}
		return &BrTable{
			Imm: &BrTableImm{
				TargetTable:   targets,
				DefaultTarget: uint32(def),
			},
		}, nil
	case RETURN:
		return &Return{}, nil
	case CALL:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(call) decode: %w", err)
		}
		return &Call{Imm: uint32(imm)}, nil
	case CALL_INDIRECT:
		index, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(call_indirect) decode: %w", err)
		}
		r, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(call_indirect) decode: %w", err)
		}
		reserved := false
		if r == 1 {
			reserved = true
		}
		return &CallIndirect{
			Imm: &CallIndirectImm{
				TypeIndex: uint32(index),
				reserved:  reserved,
			},
		}, nil
	case DROP:
		return &Drop{}, nil
	case SELECT:
		return &Select{}, nil
	case GET_LOCAL:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(get_local) decode: %w", err)
		}
		return &GetLocal{Imm: uint32(imm)}, nil
	case SET_LOCAL:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(set_local) decode: %w", err)
		}
		return &SetLocal{Imm: uint32(imm)}, nil
	case TEE_LOCAL:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(tee_local) decode: %w", err)
		}
		return &TeeLocal{Imm: uint32(imm)}, nil
	case GET_GLOBAL:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(get_global) decode: %w", err)
		}
		return &GetGlobal{Imm: uint32(imm)}, nil
	case SET_GLOBAL:
		imm, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(set_global) decode: %w", err)
		}
		return &SetGlobal{Imm: uint32(imm)}, nil
	// case I32_LOAD:
	// case I64_LOAD:
	// case F32_LOAD:
	// case F64_LOAD:
	// case I32_LOAD8_S:
	// case I32_LOAD8_U:
	// case I32_LOAD16_S:
	// case I32_LOAD16_U:
	// case I64_LOAD8_S:
	// case I64_LOAD8_U:
	// case I64_LOAD16_S:
	// case I64_LOAD16_U:
	// case I64_LOAD32_S:
	// case I64_LOAD32_U:
	// case I32_STORE:
	// case I64_STORE:
	// case F32_STORE:
	// case F64_STORE:
	// case I32_STORE8:
	// case I32_STORE16:
	// case I64_STORE8:
	// case I64_STORE16:
	// case I64_STORE32:
	// case CURRENT_MEMORY:
	// case GROW_MEMORY:
	case I32_CONST:
		imm, _, err := types.DecodeVarInt32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32_const) decode: %w", err)
		}
		return &I32Const{Imm: int32(imm)}, nil
	case I64_CONST:
		imm, _, err := types.DecodeVarInt64(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64_const) decode: %w", err)
		}
		return &I64Const{Imm: int64(imm)}, nil
	case F32_CONST:
		imm, err := types.DecodeUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32_const) decode: %w", err)
		}
		return &F32Const{Imm: imm}, nil
	case F64_CONST:
		imm, err := types.DecodeUint64(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32_const) decode: %w", err)
		}
		return &F64Const{Imm: imm}, nil
	// case I32_EQZ:
	// case I32_EQ:
	// case I32_NE:
	// case I32_LT_S:
	// case I32_LT_U:
	// case I32_GT_S:
	// case I32_GT_U:
	// case I32_LE_S:
	// case I32_LE_U:
	// case I32_GE_S:
	// case I32_GE_U:
	// case I64_EQZ:
	// case I64_EQ:
	// case I64_NE:
	// case I64_LT_S:
	// case I64_LT_U:
	// case I64_GT_S:
	// case I64_GT_U:
	// case I64_LE_S:
	// case I64_LE_U:
	// case I64_GE_S:
	// case I64_GE_U:
	// case F32_EQ:
	// case F32_NE:
	// case F32_LT:
	// case F32_GT:
	// case F32_LE:
	// case F32_GE:
	// case F64_EQ:
	// case F64_NE:
	// case F64_LT:
	// case F64_GT:
	// case F64_LE:
	// case F64_GE:
	// case I32_CLZ:
	// case I32_CTZ:
	// case I32_POPCNT:
	case I32_ADD:
		return &I32Add{}, nil
	case I32_SUB:
		return &I32Sub{}, nil
	case I32_MUL:
		return &I32Mul{}, nil
	// case I32_DIV_S:
	// case I32_DIV_U:
	// case I32_REM_S:
	// case I32_REM_U:
	// case I32_AND:
	// case I32_OR:
	// case I32_XOR:
	// case I32_SHL:
	// case I32_SHR_S:
	// case I32_SHR_U:
	// case I32_ROTL:
	// case I32_ROTR:
	// case I64_CLZ:
	// case I64_CTZ:
	// case I64_POPCNT:
	// case I64_ADD:
	// case I64_SUB:
	// case I64_MUL:
	// case I64_DIV_S:
	// case I64_DIV_U:
	// case I64_REM_S:
	// case I64_REM_U:
	// case I64_AND:
	// case I64_OR:
	// case I64_XOR:
	// case I64_SHL:
	// case I64_SHR_S:
	// case I64_SHR_U:
	// case I64_ROTL:
	// case I64_ROTR:
	// case F32_ABS:
	// case F32_NEG:
	// case F32_CEIL:
	// case F32_FLOOR:
	// case F32_TRUNC:
	// case F32_NEAREST:
	// case F32_SQRT:
	// case F32_ADD:
	// case F32_SUB:
	// case F32_MUL:
	// case F32_DIV:
	// case F32_MIN:
	// case F32_MAX:
	// case F32_COPYSIGN:
	// case F64_ABS:
	// case F64_NEG:
	// case F64_CEIL:
	// case F64_FLOOR:
	// case F64_TRUNC:
	// case F64_NEAREST:
	// case F64_SQRT:
	// case F64_ADD:
	// case F64_SUB:
	// case F64_MUL:
	// case F64_DIV:
	// case F64_MIN:
	// case F64_MAX:
	// case F64_COPYSIGN:
	// case I32_WRAP_I64:
	// case I32_TRUNC_S_F32:
	// case I32_TRUNC_U_F32:
	// case I32_TRUNC_S_F64:
	// case I32_TRUNC_U_F64:
	// case I64_EXTEND_S_I32:
	// case I64_EXTEND_U_I32:
	// case I64_TRUNC_S_F32:
	// case I64_TRUNC_U_F32:
	// case I64_TRUNC_S_F64:
	// case I64_TRUNC_U_F64:
	// case F32_CONVERT_S_I32:
	// case F32_CONVERT_U_I32:
	// case F32_CONVERT_S_I64:
	// case F32_CONVERT_U_I64:
	// case F32_DEMOTE_F64:
	// case F64_CONVERT_S_I32:
	// case F64_CONVERT_U_I32:
	// case F64_CONVERT_S_I64:
	// case F64_CONVERT_U_I64:
	// case F64_PROMOTE_F32:
	// case TRUNC_SAT:
	// case I32_REINTERPRET_F32:
	// case I64_REINTERPRET_F64:
	// case F32_REINTERPRET_I32:
	// case F64_REINTERPRET_I64:
	default:
		return nil, fmt.Errorf("%w: %x", NotImplemented, opcode)
	}
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
