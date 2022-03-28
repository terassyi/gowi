package runtime

import (
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

var (
	ExecutionErrorTypeNotMatched       error = errors.New("Execution error: type doesn't match")
	ExcetutionErrorNotConstInstruction error = errors.New("Execution error: not const instr")
	ExecutionErrorNotLocalInstruction  error = errors.New("Execution error: not local isntr")
	ExecutionErrorNotAddInstruction    error = errors.New("Execution error: not add instr")
	ExecutionErrorLocalNotExist        error = errors.New("Execution error: local values is not exist")
	ExecutionErrorArgumentTypeNotMatch error = errors.New("Execution error: argument values is not matched")
	Trap                               error = errors.New("trap")
	TrapUnreachable                    error = errors.New("trap: unreachable")
)

func (i *interpreter) labelEnd(instr instruction.Instruction) error {
	if _, err := i.stack.Frame.Pop(); err != nil {
		return fmt.Errorf("label end: %w", err)
	}
	if _, err := i.stack.Label.Pop(); err != nil {
		return fmt.Errorf("label end: %w", err)
	}
	return i.cur.update(i.stack)
}

func (i *interpreter) execDrop(instr instruction.Instruction) error {
	if _, err := i.stack.Value.Pop(); err != nil {
		return fmt.Errorf("drop: %w", err)
	}
	return nil
}

func (i *interpreter) execNop(instr instruction.Instruction) error {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-nop
	return nil
}

func (i *interpreter) execUnreachable(instr instruction.Instruction) error {
	return TrapUnreachable
}

func (i *interpreter) execSelect(instruction.Instruction) error {
	if err := i.stack.Value.Validate([]types.ValueType{types.I32}); err != nil {
		return fmt.Errorf("select :%w", err)
	}
	c, err := i.stack.Value.Pop()
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}
	val2, err := i.stack.Value.Pop()
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}
	val1, err := i.stack.Value.Pop()
	if err != nil {
		return fmt.Errorf("select: %w", err)
	}
	if instance.GetVal[value.I32](c) != value.I32(0) {
		if err := i.stack.Value.Push(val1); err != nil {
			return fmt.Errorf("select: %w", err)
		}
	} else {
		if err := i.stack.Value.Push(val2); err != nil {
			return fmt.Errorf("select: %w", err)
		}
	}
	return nil
}

func (i *interpreter) execCall(instr instruction.Instruction) error {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-call-x
	index := instruction.Imm[uint32](instr)
	fmt.Printf("function call index = %d", index)
	return nil
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
		return i.getLocal(instruction.Imm[uint32](instr), frame)
	case instruction.SET_LOCAL:
		return i.setLocal(instruction.Imm[uint32](instr), frame)
	case instruction.TEE_LOCAL:
		return i.teeLocal(instruction.Imm[uint32](instr), frame)
	default:
		return ExecutionErrorNotLocalInstruction
	}
}

func (i *interpreter) getLocal(index uint32, frame *stack.Frame) error {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-variable-mathsf-local-get-x
	if int(index) >= len(frame.Locals) {
		return ExecutionErrorLocalNotExist
	}
	if err := i.stack.Value.Push(frame.Locals[index]); err != nil {
		return err
	}
	return nil
}

func (i *interpreter) setLocal(index uint32, frame *stack.Frame) error {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-variable-mathsf-local-set-x
	if int(index) >= len(frame.Locals) {
		return ExecutionErrorLocalNotExist
	}
	val, err := i.stack.Value.Pop()
	if err != nil {
		return fmt.Errorf("set_local: %w", err)
	}
	frame.Locals[index] = val
	return nil
}

func (i *interpreter) teeLocal(index uint32, frame *stack.Frame) error {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-variable-mathsf-local-set-x
	if int(index) >= len(frame.Locals) {
		return ExecutionErrorLocalNotExist
	}
	val, err := i.stack.Value.Pop()
	if err != nil {
		return fmt.Errorf("tee_local: %w", err)
	}
	if err := i.stack.Value.Push(val); err != nil {
		return fmt.Errorf("tee_local: %w", err)
	}
	if err := i.stack.Value.Push(val); err != nil {
		return fmt.Errorf("tee_local: %w", err)
	}
	return i.setLocal(index, frame)
}

func (i *interpreter) execUnop(isntr instruction.Instruction) error {
	return nil
}

type unopFunc func(value.Number) (value.Number, error)

func (i *interpreter) unop(valType value.NumberType, f unopFunc) error {
	v, err := i.stack.Value.Top()
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	if v.ValType() != value.ValTypeNum {
		return fmt.Errorf("binop: %w", ExecutionErrorTypeNotMatched)
	}
	if v.(value.Number).NumType() != valType {
		return fmt.Errorf("binop: %w", ExecutionErrorArgumentTypeNotMatch)
	}
	val, err := i.stack.Value.Pop()
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	res, err := f(val.(value.Number))
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	if err := i.stack.Value.Push(res.ToValue()); err != nil {
		return fmt.Errorf("binop: %w", err)
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
