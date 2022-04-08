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
	ExecutionErrorDivideByZero         error = errors.New("Execution error: divide by zero")
	ExecutionErrorParse                error = errors.New("Execution error: failed to parse")
	ExecutionErrorOperation            error = errors.New("Execution error: operation error")
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
	val, err := i.stack.PopValue()
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

func (i *interpreter) execBr(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-br-l
	labelIndex := instruction.Imm[uint32](instr)
	fmt.Printf("%s label=%d\n", instr, labelIndex)
	if i.stack.Label.Len() <= int(labelIndex) {
		return instructionResultTrap, fmt.Errorf("br: the label stack must contain at least %d labels", labelIndex+1)
	}
	l, err := i.stack.Label.Ref(int(labelIndex))
	if err != nil {
		return instructionResultTrap, fmt.Errorf("br: %w", err)
	}
	values, err := i.stack.PopValuesRev(int(l.N))
	if err != nil {
		return instructionResultTrap, fmt.Errorf("br: %w", err)
	}
	for j := 0; j < int(labelIndex)+1; j++ {
		label, err := i.stack.Label.Top()
		if err != nil {
			return instructionResultTrap, fmt.Errorf("br: %w", err)
		}
		if _, err := i.stack.PopValues(int(label.ValCounter)); err != nil {
			return instructionResultTrap, fmt.Errorf("br: %w", err)
		}
		if _, err := i.stack.Label.Pop(); err != nil {
			return instructionResultTrap, fmt.Errorf("br: %w", err)
		}
	}
	for _, v := range values {
		if err := i.stack.PushValue(v); err != nil {
			return instructionResultTrap, fmt.Errorf("br: %w", err)
		}
	}
	if err := i.cur.update(i.stack); err != nil {
		return instructionResultTrap, fmt.Errorf("br: %w", err)
	}
	return instructionResultLabelEnd, nil
}

func (i *interpreter) execBrIf(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-br-l
	fmt.Printf("%s label=%x\n", instr, instruction.Imm[uint32](instr))
	if err := i.stack.Value.Validate([]types.ValueType{types.I32}); err != nil {
		return instructionResultTrap, fmt.Errorf("br_if: %w", err)
	}
	cond, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("br_if: %w", err)
	}
	switch cond {
	case value.I32(0):
		return instructionResultRunNext, nil
	default:
		return i.execBr(instr)
	}
}

func (i *interpreter) execReturn(instr instruction.Instruction) (instructionResult, error) {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-control-mathsf-return
	fmt.Println(instr)
	if i.stack.Value.Len() < int(i.cur.label.N) {
		return instructionResultTrap, fmt.Errorf("return: the value stack must have at least %d values", i.cur.label.N)
	}
	if i.stack.Frame.IsEmpty() {
		return instructionResultTrap, fmt.Errorf("return: the frame stack must have at least one vlaue")
	}
	if err := i.restoreStack(); err != nil {
		return instructionResultTrap, fmt.Errorf("label end: %w", err)
	}
	if err := i.cur.update(i.stack); err != nil {
		return instructionResultTrap, fmt.Errorf("label end: %w", err)
	}
	return instructionResultLabelEnd, nil
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
	if _, err := i.stack.PopValue(); err != nil {
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

func (i *interpreter) execSelect(instr instruction.Instruction) (instructionResult, error) {
	fmt.Println(instr)
	if err := i.stack.Value.Validate([]types.ValueType{types.I32}); err != nil {
		return instructionResultTrap, fmt.Errorf("select :%w", err)
	}
	c, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("select: %w", err)
	}
	val2, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("select: %w", err)
	}
	val1, err := i.stack.PopValue()
	if err != nil {
		return instructionResultTrap, fmt.Errorf("select: %w", err)
	}
	if instance.GetVal[value.I32](c) != value.I32(0) {
		if err := i.stack.PushValue(val1); err != nil {
			return instructionResultTrap, fmt.Errorf("select: %w", err)
		}
	} else {
		if err := i.stack.PushValue(val2); err != nil {
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
	switch instr.Opcode() {
	case instruction.I32_CONST:
		imm := instruction.Imm[int32](instr)
		fmt.Printf("%s(%v)\n", instr, imm)
		if err := i.stack.PushValue(value.I32(imm)); err != nil {
			return instructionResultTrap, fmt.Errorf("const: %w", err)
		}
		return instructionResultRunNext, nil
	case instruction.I64_CONST:
		imm := instruction.Imm[int64](instr)
		fmt.Printf("%s(%v)\n", instr, imm)
		if err := i.stack.PushValue(value.I64(imm)); err != nil {
			return instructionResultTrap, fmt.Errorf("const: %w", err)
		}
		return instructionResultRunNext, nil
	case instruction.F32_CONST:
		imm := instruction.Imm[uint32](instr)
		fmt.Printf("%s(%v)\n", instr, imm)
		if err := i.stack.PushValue(value.F32(imm)); err != nil {
			return instructionResultTrap, fmt.Errorf("const: %w", err)
		}
		return instructionResultRunNext, nil
	case instruction.F64_CONST:
		imm := instruction.Imm[uint64](instr)
		fmt.Printf("%s(%v)\n", instr, imm)
		if err := i.stack.PushValue(value.F64(imm)); err != nil {
			return instructionResultTrap, fmt.Errorf("const: %w", err)
		}
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
	if err := i.stack.PushValue(frame.Locals[index]); err != nil {
		return err
	}
	return nil
}

func (i *interpreter) setLocal(index uint32, frame *stack.Frame) error {
	// https://webassembly.github.io/spec/core/exec/instructions.html#xref-syntax-instructions-syntax-instr-variable-mathsf-local-set-x
	if int(index) >= len(frame.Locals) {
		return ExecutionErrorLocalNotExist
	}
	val, err := i.stack.PopValue()
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
	val, err := i.stack.PopValue()
	if err != nil {
		return fmt.Errorf("tee_local: %w", err)
	}
	if err := i.stack.PushValue(val); err != nil {
		return fmt.Errorf("tee_local: %w", err)
	}
	if err := i.stack.PushValue(val); err != nil {
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
	val, err := i.stack.PopValue()
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	res, err := f(val.(value.Number))
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	if err := i.stack.PushValue(res.ToValue()); err != nil {
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
		if err := i.binop(value.NumTypeI64, add); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_SUB:
		if err := i.binop(value.NumTypeI32, sub); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_SUB:
		if err := i.binop(value.NumTypeI64, sub); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_MUL:
		if err := i.binop(value.NumTypeI32, mul); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_MUL:
		if err := i.binop(value.NumTypeI64, mul); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_DIV_S:
		if err := i.binop(value.NumTypeI32, divs); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_DIV_S:
		if err := i.binop(value.NumTypeI64, divs); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_DIV_U:
		if err := i.binop(value.NumTypeI32, divu); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_DIV_U:
		if err := i.binop(value.NumTypeI64, divu); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_REM_S:
		if err := i.binop(value.NumTypeI32, rems); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_REM_S:
		if err := i.binop(value.NumTypeI64, rems); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_REM_U:
		if err := i.binop(value.NumTypeI32, remu); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_REM_U:
		if err := i.binop(value.NumTypeI64, remu); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_AND:
		if err := i.binop(value.NumTypeI32, and); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_AND:
		if err := i.binop(value.NumTypeI64, and); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_OR:
		if err := i.binop(value.NumTypeI32, or); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_OR:
		if err := i.binop(value.NumTypeI64, or); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I32_XOR:
		if err := i.binop(value.NumTypeI32, xor); err != nil {
			return instructionResultTrap, err
		}
	case instruction.I64_XOR:
		if err := i.binop(value.NumTypeI64, xor); err != nil {
			return instructionResultTrap, err
		}
	case instruction.F32_ADD:
	case instruction.F64_ADD:
	default:
		return instructionResultTrap, instruction.NotImplemented
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
	values, err := i.stack.PopValuesRev(2)
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	res, err := f(values[0].(value.Number), values[1].(value.Number))
	if err != nil {
		return fmt.Errorf("binop: %w", err)
	}
	if err := i.stack.PushValue(res.ToValue()); err != nil {
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
	case value.NumTypeI64:
		return value.I64(int64(value.GetNum[value.I64](a)) + int64(value.GetNum[value.I64](b))), nil
	}
	return nil, nil
}

func sub(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		return value.I32(int32(value.GetNum[value.I32](a)) - int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		return value.I64(int64(value.GetNum[value.I64](a)) - int64(value.GetNum[value.I64](b))), nil
	}
	return nil, nil
}

func mul(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		return value.I32(int32(value.GetNum[value.I32](a)) * int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		return value.I64(int64(value.GetNum[value.I64](a)) * int64(value.GetNum[value.I64](b))), nil
	}
	return nil, nil
}

func divs(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		if value.GetNum[value.I32](b) == value.I32(0) {
			return nil, ExecutionErrorDivideByZero
		}
		return value.I32(int32(value.GetNum[value.I32](a)) / int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		if value.GetNum[value.I64](b) == value.I64(0) {
			return nil, ExecutionErrorDivideByZero
		}
		return value.I64(int64(value.GetNum[value.I64](a)) / int64(value.GetNum[value.I64](b))), nil
	}
	return nil, nil
}

func divu(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		if value.GetNum[value.I32](b) == value.I32(0) {
			return nil, ExecutionErrorDivideByZero
		}
		ua, err := value.GetNum[value.I32](a).ToUint32()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		ub, err := value.GetNum[value.I32](b).ToUint32()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		return value.I32(ua / ub), nil
	case value.NumTypeI64:
		if value.GetNum[value.I64](b) == value.I64(0) {
			return nil, ExecutionErrorDivideByZero
		}
		ua, err := value.GetNum[value.I64](a).ToUint64()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		ub, err := value.GetNum[value.I64](b).ToUint64()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		return value.I64(ua / ub), nil
	}
	return nil, nil
}

func rems(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		if value.GetNum[value.I32](b) == value.I32(0) {
			return nil, ExecutionErrorDivideByZero
		}
		return value.I32(int32(value.GetNum[value.I32](a)) % int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		if value.GetNum[value.I64](b) == value.I64(0) {
			return nil, ExecutionErrorDivideByZero
		}
		return value.I64(int64(value.GetNum[value.I64](a)) % int64(value.GetNum[value.I64](b))), nil
	}
	return nil, nil
}

func remu(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		if value.GetNum[value.I32](b) == value.I32(0) {
			return nil, ExecutionErrorDivideByZero
		}
		ua, err := value.GetNum[value.I32](a).ToUint32()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		ub, err := value.GetNum[value.I32](b).ToUint32()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		div, err := divu(a, b)
		if err != nil {
			return nil, ExecutionErrorOperation
		}
		udiv, err := value.GetNum[value.I32](div).ToUint32()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		res := ua - ub*udiv
		return value.I32(res), nil

	case value.NumTypeI64:
		if value.GetNum[value.I64](b) == value.I64(0) {
			return nil, ExecutionErrorDivideByZero
		}
		ua, err := value.GetNum[value.I64](a).ToUint64()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		ub, err := value.GetNum[value.I64](b).ToUint64()
		if err != nil {
			return nil, ExecutionErrorParse
		}
		return value.I64(ua % ub), nil
	}
	return nil, nil
}

func and(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		return value.I32(int32(value.GetNum[value.I32](a)) & int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		return value.I64(int64(value.GetNum[value.I64](a)) & int64(value.GetNum[value.I64](b))), nil
	default:
		return nil, ExecutionErrorOperation
	}
}

func or(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		return value.I32(int32(value.GetNum[value.I32](a)) | int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		return value.I64(int64(value.GetNum[value.I64](a)) | int64(value.GetNum[value.I64](b))), nil
	default:
		return nil, ExecutionErrorOperation
	}
}

func xor(a, b value.Number) (value.Number, error) {
	if a.NumType() != b.NumType() {
		return nil, ExecutionErrorArgumentTypeNotMatch
	}
	switch a.NumType() {
	case value.NumTypeI32:
		return value.I32(int32(value.GetNum[value.I32](a)) ^ int32(value.GetNum[value.I32](b))), nil
	case value.NumTypeI64:
		return value.I64(int64(value.GetNum[value.I64](a)) ^ int64(value.GetNum[value.I64](b))), nil
	default:
		return nil, ExecutionErrorOperation
	}
}
