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

type instructionResult uint8

const (
	instructionResultRunNext    instructionResult = iota
	instructionResultCallFunc   instructionResult = iota
	instructionResultEnterBlock instructionResult = iota
	instructionResultLabelEnd   instructionResult = iota
	instructionResultReturn     instructionResult = iota
	instructionResultTrap       instructionResult = iota
)

type blockStack struct {
	inner []block
}

func newBlockStack() *blockStack {
	return &blockStack{
		inner: make([]block, 0, 1024),
	}
}

func (b *blockStack) push(op instruction.Opcode) error {
	switch op {
	case instruction.BLOCK:
		b.inner = append(b.inner, &blockBlock{})
	case instruction.LOOP:
		b.inner = append(b.inner, &blockLoop{})
	case instruction.IF:
		b.inner = append(b.inner, &blockIf{})
	default:
		return fmt.Errorf("Instruction is not structured instruction")
	}
	return nil
}

func (b *blockStack) pop() (block, error) {
	if b.len() == 0 {
		return nil, fmt.Errorf("block stack is empty")
	}
	val := b.inner[len(b.inner)-1]
	b.inner = b.inner[:len(b.inner)-1]
	return val, nil
}

func (b *blockStack) top() block {
	if b.isEmpty() {
		return nil
	}
	return b.inner[len(b.inner)-1]
}

func (b *blockStack) len() int {
	return len(b.inner)
}

func (b *blockStack) isEmpty() bool {
	if b.len() == 0 {
		return true
	}
	return false
}

type blockType uint8

const (
	blockTypeBlock blockType = iota
	blockTypeLoop  blockType = iota
	blockTypeIf    blockType = iota
)

type block interface {
	typ() blockType
}

type blockBlock struct{}

type blockLoop struct{}

type blockIf struct {
	els bool
}

func (*blockBlock) typ() blockType {
	return blockTypeBlock
}

func (*blockLoop) typ() blockType {
	return blockTypeLoop
}

func (*blockIf) typ() blockType {
	return blockTypeIf
}

func (i *interpreter) execBlock(instr instruction.Instruction) (instructionResult, error) {
	fmt.Printf("%s: %x\n", instr, instruction.Imm[types.BlockType](instr))
	funcType, err := i.expand(instruction.Imm[types.BlockType](instr))
	if err != nil {
		return instructionResultTrap, fmt.Errorf("block: %w", err)
	}
	label, l, err := i.labelBlock(funcType)
	if err != nil {
		return instructionResultTrap, fmt.Errorf("block: %w", err)
	}
	i.cur.label.Sp += l
	if err := i.stack.Label.Push(*label); err != nil {
		return instructionResultTrap, fmt.Errorf("block: %w", err)
	}
	return instructionResultEnterBlock, nil
}

func (i *interpreter) execLoop(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-loop-xref-syntax-instructions-syntax-blocktype-mathit-blocktype-xref-syntax-instructions-syntax-instr-mathit-instr-ast-xref-syntax-instructions-syntax-instr-control-mathsf-end
	fmt.Printf("%s: %x\n", instr, instruction.Imm[types.BlockType](instr))
	funcType, err := i.expand(instruction.Imm[types.BlockType](instr))
	if err != nil {
		return instructionResultTrap, fmt.Errorf("loop: %w", err)
	}
	label, l, err := i.labelBlock(funcType)
	if err != nil {
		return instructionResultTrap, fmt.Errorf("loop: %w", err)
	}
	i.cur.label.Sp += l
	if err := i.stack.Label.Push(*label); err != nil {
		return instructionResultTrap, fmt.Errorf("loop: %w", err)
	}
	return instructionResultEnterBlock, nil
}

func (i *interpreter) execIf(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-if-xref-syntax-instructions-syntax-blocktype-mathit-blocktype-xref-syntax-instructions-syntax-instr-mathit-instr-1-ast-xref-syntax-instructions-syntax-instr-control-mathsf-else-xref-syntax-instructions-syntax-instr-mathit-instr-2-ast-xref-syntax-instructions-syntax-instr-control-mathsf-end
	fmt.Printf("%s: %x\n", instr, instruction.Imm[types.BlockType](instr))
	funcType, err := i.expand(instruction.Imm[types.BlockType](instr))
	if err != nil {
		return instructionResultTrap, fmt.Errorf("if: %w", err)
	}
	label, l, err := i.labelBlock(funcType)
	if err := i.stack.Value.Validate([]types.ValueType{types.I32}); err != nil {
		return instructionResultTrap, fmt.Errorf("if: %w", err)
	}
	fmt.Printf("[debug] stack frame pointer pop len: %d\n", l)
	i.cur.label.Sp += l
	val, err := i.stack.Value.Pop()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("if: %w", err)
	}
	if val.(value.Number).NumType() != value.NumTypeI32 {
		return instructionResultTrap, fmt.Errorf("if: i32 is expected, got %s", val.(value.Number))
	}
	if i.stack.Value.Len() < len(funcType.Params) {
		return instructionResultTrap, fmt.Errorf("if: %d values is requied on the value stack", len(funcType.Params))
	}
	// branch
	condLabel, err := ifElseLabel(label, instance.GetVal[value.I32](val))
	if err != nil {
		return instructionResultTrap, fmt.Errorf("if: %w", err)
	}
	// push to label stack
	if err := i.stack.Label.Push(*condLabel); err != nil {
		return instructionResultTrap, fmt.Errorf("loop: %w", err)
	}
	return instructionResultEnterBlock, nil
}

func (i *interpreter) labelBlock(funcType *types.FuncType) (*stack.Label, int, error) {
	instrs := make([]instruction.Instruction, 0)
	bs := newBlockStack()
	bs.push(i.cur.label.Instructions[i.cur.label.Sp].Opcode())
	for sp := i.cur.label.Sp + 1; sp < len(i.cur.label.Instructions); sp++ {
		instrs = append(instrs, i.cur.label.Instructions[sp])
		if i.cur.label.Instructions[sp].Opcode() == instruction.BLOCK ||
			i.cur.label.Instructions[sp].Opcode() == instruction.LOOP ||
			i.cur.label.Instructions[sp].Opcode() == instruction.IF {
			bs.push(i.cur.label.Instructions[sp].Opcode())
		}
		if i.cur.label.Instructions[sp].Opcode() == instruction.END {
			_, err := bs.pop()
			if err != nil {
				return nil, 0, fmt.Errorf("label block: %w", err)
			}
			if bs.isEmpty() {
				break
			}
		}
	}
	return &stack.Label{Instructions: instrs, N: uint8(len(funcType.Returns)), Sp: 0, Flag: false}, len(instrs), nil
}

func ifElseLabel(label *stack.Label, cond value.I32) (*stack.Label, error) {
	ifBlock := make([]instruction.Instruction, 0, len(label.Instructions))
	elseBlock := make([]instruction.Instruction, 0, len(label.Instructions))
	bs := newBlockStack()
	bs.push(instruction.IF)
	splitIndex := -1
	for j, instr := range label.Instructions {
		if instr.Opcode() == instruction.BLOCK ||
			instr.Opcode() == instruction.LOOP ||
			instr.Opcode() == instruction.IF {
			if err := bs.push(instr.Opcode()); err != nil {
				return nil, fmt.Errorf("if_else label: %w", err)
			}
		}
		if instr.Opcode() == instruction.ELSE {
			if bs.top() == nil || bs.top().typ() != blockTypeIf {
				return nil, fmt.Errorf("if_else label: invalie else")
			}
			if bs.len() == 1 {
				splitIndex = j
			}
		}
		if instr.Opcode() == instruction.END {
			_, err := bs.pop()
			if err != nil {
				return nil, fmt.Errorf("if_else label: %w", err)
			}
			if bs.isEmpty() {
				break
			}
		}
	}
	if splitIndex != -1 {
		ifBlock = append(ifBlock, label.Instructions[:splitIndex]...)
		ifBlock = append(ifBlock, &instruction.End{})
		elseBlock = append(elseBlock, label.Instructions[splitIndex+1:]...)
		if len(elseBlock) == 0 {
			elseBlock = append(elseBlock, &instruction.End{})
		}
	} else {
		ifBlock = append(ifBlock, label.Instructions...)
		elseBlock = append(elseBlock, &instruction.End{})
	}
	condLabel := &stack.Label{N: label.N, Sp: label.Sp, Flag: false}
	if cond != value.I32(0) {
		condLabel.Instructions = ifBlock
	} else {
		condLabel.Instructions = elseBlock
	}
	return condLabel, nil
}

func (i *interpreter) execLabelEnd(instr instruction.Instruction) (instructionResult, error) {
	fmt.Println(instr.String())
	if err := i.restoreStack(); err != nil {
		return instructionResultTrap, fmt.Errorf("label end: %w", err)
	}
	if err := i.cur.update(i.stack); err != nil {
		return instructionResultTrap, fmt.Errorf("label end: %w", err)
	}
	return instructionResultLabelEnd, nil
}

func (i *interpreter) execDrop(instr instruction.Instruction) (instructionResult, error) {
	if _, err := i.stack.Value.Pop(); err != nil {
		return instructionResultTrap, fmt.Errorf("drop: %w", err)
	}
	return instructionResultRunNext, nil
}

func (i *interpreter) execNop(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-nop
	fmt.Println(instr.String())
	return instructionResultRunNext, nil
}

func (i *interpreter) execUnreachable(instr instruction.Instruction) (instructionResult, error) {
	return instructionResultTrap, TrapUnreachable
}

func (i *interpreter) execSelect(instruction.Instruction) (instructionResult, error) {
	if err := i.stack.Value.Validate([]types.ValueType{types.I32}); err != nil {
		return instructionResultTrap, fmt.Errorf("select :%w", err)
	}
	c, err := i.stack.Value.Pop()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("select: %w", err)
	}
	val2, err := i.stack.Value.Pop()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("select: %w", err)
	}
	val1, err := i.stack.Value.Pop()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("select: %w", err)
	}
	if instance.GetVal[value.I32](c) != value.I32(0) {
		if err := i.stack.Value.Push(val1); err != nil {
			return instructionResultTrap, fmt.Errorf("select: %w", err)
		}
	} else {
		if err := i.stack.Value.Push(val2); err != nil {
			return instructionResultTrap, fmt.Errorf("select: %w", err)
		}
	}
	return instructionResultRunNext, nil
}

func (i *interpreter) execCall(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-call-x
	index := instruction.Imm[uint32](instr)
	fmt.Printf("function call index = %d\n", index)
	i.f = i.cur.frame.Module.FuncAddrs[index]
	// return i.invokeFunction(f)
	return instructionResultCallFunc, nil
}

func (i *interpreter) execConst(instr instruction.Instruction) (instructionResult, error) {
	fmt.Printf("%s(%v)\n", instr, instruction.Imm[int32](instr))
	switch instr.Opcode() {
	case instruction.I32_CONST:
		imm := instruction.Imm[int32](instr)
		if err := i.stack.Value.Push(value.I32(imm)); err != nil {
			return instructionResultTrap, fmt.Errorf("const: %w", err)
		}
		return instructionResultRunNext, nil
	case instruction.I64_CONST:
		return instructionResultRunNext, nil
	case instruction.F32_CONST:
		return instructionResultRunNext, nil
	case instruction.F64_CONST:
		return instructionResultRunNext, nil
	default:
		return instructionResultTrap, ExcetutionErrorNotConstInstruction
	}
}

func (i *interpreter) execLocal(instr instruction.Instruction, frame *stack.Frame) (instructionResult, error) {
	fmt.Printf("%s(%v)\n", instr, instruction.Imm[uint32](instr))
	switch instr.Opcode() {
	case instruction.GET_LOCAL:
		if err := i.getLocal(instruction.Imm[uint32](instr), frame); err != nil {
			return instructionResultTrap, fmt.Errorf("get_local: %w", err)
		}
	case instruction.SET_LOCAL:
		if err := i.setLocal(instruction.Imm[uint32](instr), frame); err != nil {
			return instructionResultTrap, fmt.Errorf("set_local: %w", err)
		}
	case instruction.TEE_LOCAL:
		if err := i.teeLocal(instruction.Imm[uint32](instr), frame); err != nil {
			return instructionResultTrap, fmt.Errorf("tee_local: %w", err)
		}
	default:
		return instructionResultTrap, ExecutionErrorNotLocalInstruction
	}
	return instructionResultRunNext, nil
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

func (i *interpreter) execBinop(instr instruction.Instruction) (instructionResult, error) {
	fmt.Println(instr.String())
	switch instr.Opcode() {
	case instruction.I32_ADD:
		if err := i.binop(value.NumTypeI32, add); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_ADD:
	case instruction.F32_ADD:
	case instruction.F64_ADD:
	default:
		return instructionResultTrap, ExecutionErrorNotAddInstruction
	}
	return instructionResultRunNext, nil
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
