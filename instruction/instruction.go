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
	String() string
	imm() any
	ImmString() string
}

type None struct{}

var NoImm None = None{}

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
			Imm: BrTableImm{
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
			Imm: CallIndirectImm{
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
	case I32_LOAD:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.load): %w", err)
		}
		return &I32Load{Imm: *imm}, nil
	case I64_LOAD:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load): %w", err)
		}
		return &I64Load{Imm: *imm}, nil
	// case F32_LOAD:
	// case F64_LOAD:
	case I32_LOAD8_S:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.load8_s): %w", err)
		}
		return &I32Load8S{Imm: *imm}, nil
	case I32_LOAD8_U:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.load8_u): %w", err)
		}
		return &I32Load8U{Imm: *imm}, nil
	case I32_LOAD16_S:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.load16_s): %w", err)
		}
		return &I32Load16S{Imm: *imm}, nil
	case I32_LOAD16_U:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.load16_u): %w", err)
		}
		return &I32Load16U{Imm: *imm}, nil
	case I64_LOAD8_S:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load8_s): %w", err)
		}
		return &I64Load8S{Imm: *imm}, nil
	case I64_LOAD8_U:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load8_u): %w", err)
		}
		return &I64Load8U{Imm: *imm}, nil
	case I64_LOAD16_S:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load16_s): %w", err)
		}
		return &I64Load16S{Imm: *imm}, nil
	case I64_LOAD16_U:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load16_u): %w", err)
		}
		return &I64Load16U{Imm: *imm}, nil
	case I64_LOAD32_S:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load32_s): %w", err)
		}
		return &I64Load32S{Imm: *imm}, nil
	case I64_LOAD32_U:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.load32_u): %w", err)
		}
		return &I64Load32U{Imm: *imm}, nil
	case I32_STORE:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.store): %w", err)
		}
		return &I32Store{Imm: *imm}, nil
	case I64_STORE:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.store): %w", err)
		}
		return &I64Store{Imm: *imm}, nil
	// case F32_STORE:
	// case F64_STORE:
	case I32_STORE8:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.store8): %w", err)
		}
		return &I32Store8{Imm: *imm}, nil
	case I32_STORE16:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i32.store16): %w", err)
		}
		return &I32Store16{Imm: *imm}, nil
	case I64_STORE8:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.store8): %w", err)
		}
		return &I64Store8{Imm: *imm}, nil
	case I64_STORE16:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.store16): %w", err)
		}
		return &I64Store16{Imm: *imm}, nil
	case I64_STORE32:
		imm, err := newMemImm(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(i64.store32): %w", err)
		}
		return &I64Store32{Imm: *imm}, nil
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
			return nil, fmt.Errorf("Instruction(f32_const) decode: %w", err)
		}
		return &F32Const{Imm: imm}, nil
	case F64_CONST:
		imm, err := types.DecodeUint64(buf)
		if err != nil {
			return nil, fmt.Errorf("Instruction(f64_const) decode: %w", err)
		}
		return &F64Const{Imm: imm}, nil
	case I32_EQZ:
		return &I32Eqz{}, nil
	case I32_EQ:
		return &I32Eq{}, nil
	case I32_NE:
		return &I32Ne{}, nil
	case I32_LT_S:
		return &I32LtS{}, nil
	case I32_LT_U:
		return &I32LtU{}, nil
	case I32_GT_S:
		return &I32GtS{}, nil
	case I32_GT_U:
		return &I32GtU{}, nil
	case I32_LE_S:
		return &I32LeS{}, nil
	case I32_LE_U:
		return &I32LeU{}, nil
	case I32_GE_S:
		return &I32GeS{}, nil
	case I32_GE_U:
		return &I32GeU{}, nil
	case I64_EQZ:
		return &I64Eqz{}, nil
	case I64_EQ:
		return &I64Eq{}, nil
	case I64_NE:
		return &I64Ne{}, nil
	case I64_LT_S:
		return &I64LtS{}, nil
	case I64_LT_U:
		return &I64LtU{}, nil
	case I64_GT_S:
		return &I64GtS{}, nil
	case I64_GT_U:
		return &I64GtU{}, nil
	case I64_LE_S:
		return &I64LeS{}, nil
	case I64_LE_U:
		return &I64LeU{}, nil
	case I64_GE_S:
		return &I64GeS{}, nil
	case I64_GE_U:
		return &I64GeU{}, nil
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
	case I32_CLZ:
		return &I32Clz{}, nil
	case I32_CTZ:
		return &I32Ctz{}, nil
	case I32_POPCNT:
		return &I32Popcnt{}, nil
	case I32_ADD:
		return &I32Add{}, nil
	case I32_SUB:
		return &I32Sub{}, nil
	case I32_MUL:
		return &I32Mul{}, nil
	case I32_DIV_S:
		return &I32DivS{}, nil
	case I32_DIV_U:
		return &I32DivU{}, nil
	case I32_REM_S:
		return &I32RemS{}, nil
	case I32_REM_U:
		return &I32RemU{}, nil
	case I32_AND:
		return &I32And{}, nil
	case I32_OR:
		return &I32Or{}, nil
	case I32_XOR:
		return &I32Xor{}, nil
	case I32_SHL:
		return &I32Shl{}, nil
	case I32_SHR_S:
		return &I32ShrS{}, nil
	case I32_SHR_U:
		return &I32ShrU{}, nil
	case I32_ROTL:
		return &I32RotL{}, nil
	case I32_ROTR:
		return &I32RotR{}, nil
	case I64_CLZ:
		return &I64Clz{}, nil
	case I64_CTZ:
		return &I64Ctz{}, nil
	case I64_POPCNT:
		return &I64Popcnt{}, nil
	case I64_ADD:
		return &I64Add{}, nil
	case I64_SUB:
		return &I64Sub{}, nil
	case I64_MUL:
		return &I64Mul{}, nil
	case I64_DIV_S:
		return &I64DivS{}, nil
	case I64_DIV_U:
		return &I64DivU{}, nil
	case I64_REM_S:
		return &I64RemS{}, nil
	case I64_REM_U:
		return &I64RemU{}, nil
	case I64_AND:
		return &I64And{}, nil
	case I64_OR:
		return &I64Or{}, nil
	case I64_XOR:
		return &I64Xor{}, nil
	case I64_SHL:
		return &I64Shl{}, nil
	case I64_SHR_S:
		return &I64ShrS{}, nil
	case I64_SHR_U:
		return &I64ShrU{}, nil
	case I64_ROTL:
		return &I64RotL{}, nil
	case I64_ROTR:
		return &I64RotR{}, nil
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
		return nil, fmt.Errorf("%w: 0x%x", NotImplemented, opcode)
	}
}

func Imm[T any](instr Instruction) T {
	return instr.imm().(T)
}
