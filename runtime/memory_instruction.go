package runtime

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

var (
	MemoryDoesNotHaveEnoughLength error = errors.New("memory doesn't have enough length")
)

func (i *interpreter) execLoad(instr instruction.Instruction) (instructionResult, error) {
	imm := instruction.Imm[instruction.MemoryImm](instr)
	if len(i.cur.frame.Module.MemAddrs) == 0 || i.cur.frame.Module.MemAddrs == nil {
		return instructionResultTrap, fmt.Errorf("load: memory instance is not exist")
	}
	mem := i.cur.frame.Module.MemAddrs[0]
	// if err := i.stack.Value.Validate([]types.ValueType{types.I32}); err != nil {
	if err := i.stack.ValidateValue([]types.ValueType{types.I32}); err != nil {
		return instructionResultTrap, fmt.Errorf("load: %w", err)
	}
	v, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("load: %w", err)
	}
	ea := instance.GetVal[value.I32](v).Unsigned() + imm.Offset
	switch instr.Opcode() {
	case instruction.I32_LOAD:
		if int(ea+4) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.load: memory doesn't have enough length")
		}
		val := binary.LittleEndian.Uint32(load(mem, ea, 4))
		if err := i.stack.PushValue(value.I32(val)); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.load: %w", err)
		}
	case instruction.I64_LOAD:
		if int(ea+8) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load: memory doesn't have enough length")
		}
		val := binary.LittleEndian.Uint64(load(mem, ea, 8))
		if err := i.stack.PushValue(value.I64(val)); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load: %w", err)
		}
	case instruction.I32_LOAD8_S:
		if int(ea+1) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.load8_s: %w", err)
		}
		val := bytesToVal[int8](load(mem, ea, 1))
		if err := i.stack.PushValue(value.NewI32(int32(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.load8_s: %w", err)
		}
	case instruction.I64_LOAD8_S:
		if int(ea+1) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load8_s: %w", err)
		}
		val := bytesToVal[int8](load(mem, ea, 1))
		if err := i.stack.PushValue(value.NewI64(int64(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load8_s: %w", err)
		}
	case instruction.I32_LOAD8_U:
		if int(ea+1) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.load8_u: %w", err)
		}
		val := bytesToVal[uint8](load(mem, ea, 1))
		if err := i.stack.PushValue(value.I32(uint32(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.load8_u: %w", err)
		}
	case instruction.I64_LOAD8_U:
		if int(ea+1) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load8_u: %w", err)
		}
		val := bytesToVal[uint8](load(mem, ea, 1))
		if err := i.stack.PushValue(value.I64(uint64(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load8_u: %w", err)
		}
	case instruction.I32_LOAD16_S:
		if int(ea+2) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.load16_s: %w", err)
		}
		val := bytesToVal[int16](load(mem, ea, 2))
		if err := i.stack.PushValue(value.NewI32(int32(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.load16_s: %w", err)
		}
	case instruction.I64_LOAD16_S:
		if int(ea+2) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load16_s: %w", err)
		}
		val := bytesToVal[int16](load(mem, ea, 2))
		if err := i.stack.PushValue(value.NewI64(int64(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load16_s: %w", err)
		}
	case instruction.I32_LOAD16_U:
		if int(ea+2) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.load16_u: %w", err)
		}
		val := bytesToVal[uint16](load(mem, ea, 2))
		if err := i.stack.PushValue(value.I32(uint32(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.load16_u: %w", err)
		}
	case instruction.I64_LOAD16_U:
		if int(ea+2) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load16_u: %w", err)
		}
		val := bytesToVal[uint16](load(mem, ea, 2))
		if err := i.stack.PushValue(value.I64(uint64(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load16_u: %w", err)
		}
	case instruction.I64_LOAD32_S:
		if int(ea+4) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load32_s: %w", err)
		}
		val := bytesToVal[int32](load(mem, ea, 4))
		if err := i.stack.PushValue(value.NewI64(int64(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load32_s: %w", err)
		}
	case instruction.I64_LOAD32_U:
		if int(ea+4) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.load32_u: %w", err)
		}
		val := bytesToVal[uint32](load(mem, ea, 4))
		if err := i.stack.PushValue(value.I64(uint64(val))); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.load32_u: %w", err)
		}
	default:
		return instructionResultTrap, instruction.NotImplemented
	}
	return instructionResultRunNext, nil
}

func load(mem *instance.Memory, offset, length uint32) []byte {
	return mem.Data[offset : offset+length]
}

func (i *interpreter) execStore(instr instruction.Instruction) (instructionResult, error) {
	imm := instruction.Imm[instruction.MemoryImm](instr)
	if len(i.cur.frame.Module.MemAddrs) == 0 || i.cur.frame.Module.MemAddrs == nil {
		return instructionResultTrap, fmt.Errorf("store: memory instance is not exist")
	}
	mem := i.cur.frame.Module.MemAddrs[0]
	v, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("store: %w", err)
	}
	e, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("store: %w", err)
	}
	ea := instance.GetVal[value.I32](e).Unsigned() + imm.Offset
	switch instr.Opcode() {
	case instruction.I32_STORE:
		if int(ea+4) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.store: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I32](v).Unsigned(), ea, 4); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.store: %w", err)
		}
	case instruction.I64_STORE:
		if int(ea+8) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.store: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I64](v).Unsigned(), ea, 8); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.store: %w", err)
		}
	case instruction.I32_STORE8:
		if int(ea+1) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.store8: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I32](v).Unsigned(), ea, 1); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.store8: %w", err)
		}
	case instruction.I64_STORE8:
		if int(ea+1) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.store8: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I64](v).Unsigned(), ea, 1); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.store8: %w", err)
		}
	case instruction.I32_STORE16:
		if int(ea+2) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i32.store16: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I32](v).Unsigned(), ea, 2); err != nil {
			return instructionResultTrap, fmt.Errorf("i32.store16: %w", err)
		}
	case instruction.I64_STORE16:
		if int(ea+2) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.store16: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I64](v).Unsigned(), ea, 2); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.store16: %w", err)
		}
	case instruction.I64_STORE32:
		if int(ea+4) > len(mem.Data) {
			return instructionResultTrap, fmt.Errorf("i64.store32: memory doesn't have enough length")
		}
		if err := store(mem, instance.GetVal[value.I64](v).Unsigned(), ea, 4); err != nil {
			return instructionResultTrap, fmt.Errorf("i64.store32: %w", err)
		}
	default:
		return instructionResultTrap, instruction.NotImplemented
	}
	return instructionResultRunNext, nil
}

func store[T uint32 | uint64](mem *instance.Memory, val T, offset, length uint32) error {
	buf := bytes.NewBuffer(make([]byte, 0, 8))
	if err := binary.Write(buf, binary.LittleEndian, val); err != nil {
		return err
	}
	for i, b := range buf.Bytes() {
		mem.Data[int(offset)+i] = b
	}
	return nil
}
