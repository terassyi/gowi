package runtime

import (
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
)

var (
	ExecutionErrorTypeNotMatched       error = errors.New("Execution error: type doesn't match")
	ExcetutionErrorNotConstInstruction error = errors.New("Execution error: not const instr")
	ExecutionErrorNotLocalInstruction  error = errors.New("Execution error: not local isntr")
	ExecutionErrorNotAddInstruction    error = errors.New("Execution error: not add instr")
	ExecutionErrorLocalNotExist        error = errors.New("Execution error: local values is not exist")
	ExecutionErrorArgumentTypeNotMatch error = errors.New("Execution error: argument values is not matched")
)

func (i *interpreter) labelEnd(instr instruction.Instruction) error {
	// f, err := i.stack.Frame.Pop()
	// if err != nil {
	// 	return fmt.Errorf("label end: %w", err)
	// }
	// l, err := i.stack.Label.Pop()
	// if err != nil {
	// 	return fmt.Errorf("label end: %w", err)
	// }
	if _, err := i.stack.Frame.Pop(); err != nil {
		return fmt.Errorf("label end: %w", err)
	}
	if _, err := i.stack.Label.Pop(); err != nil {
		return fmt.Errorf("label end: %w", err)
	}
	return i.cur.update(i.stack)
}

func (i *interpreter) execConst(instr instruction.Instruction) error {
	fmt.Printf("%s(%v)\n", instr, instruction.Imm[int32](instr))
	switch instr.Opcode() {
	case instruction.I32_CONST:
		imm := instruction.Imm[int32](instr)
		if err := i.stack.Value.Push(value.I32(imm)); err != nil {
			return nil
		}
		return nil
	case instruction.I64_CONST:
		return nil
	case instruction.F32_CONST:
		return nil
	case instruction.F64_CONST:
		return nil
	default:
		return ExcetutionErrorNotConstInstruction
	}
}

func (i *interpreter) execLocal(instr instruction.Instruction, frame *stack.Frame) error {
	fmt.Printf("%s(%v)\n", instr, instruction.Imm[uint32](instr))
	switch instr.Opcode() {
	case instruction.GET_LOCAL:
		index := instruction.Imm[uint32](instr)
		if int(index) >= len(frame.Locals) {
			return ExecutionErrorLocalNotExist
		}
		if err := i.stack.Value.Push(frame.Locals[index]); err != nil {
			return err
		}
	case instruction.SET_LOCAL:
	case instruction.TEE_LOCAL:
	default:
		return ExecutionErrorNotLocalInstruction
	}
	return nil
}

func (i *interpreter) execBinop(instr instruction.Instruction) error {
	fmt.Println(instr.String())
	switch instr.Opcode() {
	case instruction.I32_ADD:
		return i.binop(value.NumTypeI32, add)
	case instruction.I64_ADD:
	case instruction.F32_ADD:
	case instruction.F64_ADD:
	default:
		return ExecutionErrorNotAddInstruction
	}
	return nil
}

type binopFunc func(value.Number, value.Number) (value.Number, error)

// https://webassembly.github.io/spec/core/exec/instructions.html#t-mathsf-xref-syntax-instructions-syntax-binop-mathit-binop
func (i *interpreter) binop(valType value.NumberType, f binopFunc) error {
	vallds, err := i.stack.Value.RefNRev(2)
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	for _, v := range vallds {
		if v.ValType() != value.ValTypeNum {
			return fmt.Errorf("binop: %w", ExecutionErrorTypeNotMatched)
		}
		if v.(value.Number).NumType() != valType || v.(value.Number).NumType() != valType {
			return fmt.Errorf("binop: %w", ExecutionErrorArgumentTypeNotMatch)
		}
	}
	values, err := i.stack.Value.PopNRev(2)
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	res, err := f(values[0].(value.Number), values[1].(value.Number))
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	if err := i.stack.Value.Push(res.ToValue()); err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	return nil
}

func add(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		return value.I32(int32(value.GetNum[value.I32](a)) + int32(value.GetNum[value.I32](b))), nil
	}
	return nil, nil
}
