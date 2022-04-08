package runtime

import (
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/debugger"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
	"github.com/terassyi/gowi/validator"
)

var (
	FunctionIsRequired            error = errors.New("External value type function is required")
	FunctionParamsDoesntMatch     error = errors.New("Number of function parameters doesn't match")
	FunctionParamTypesDoesntMatch error = errors.New("Function parameter type doesn't match")
)

type Interpreter interface {
	Invoke(string, []value.Value) ([]value.Value, error)
}

type interpreter struct {
	instance *instance.Module
	stack    *stack.Stack
	debubber *debugger.Debugger
	f        *instance.Function // next function
	cur      *current
}

type current struct {
	frame *stack.Frame
	label *stack.Label
}

func (c *current) update(s *stack.Stack) error {
	if err := c.updateLabel(s); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return c.updateFrame(s)
}

func (c *current) updateLabel(s *stack.Stack) error {
	label, err := s.Label.Top()
	if err != nil {
		return err
	}
	c.label = label
	return nil
}

func (c *current) updateFrame(s *stack.Stack) error {
	frame, err := s.Frame.Top()
	if err != nil {
		return err
	}
	c.frame = frame
	return nil
}

// instanciate an interpreter
// https://webassembly.github.io/spec/core/exec/modules.html#instantiation
func New(mod *structure.Module, externalvals []instance.ExternalValue, debugLevel debugger.DebugLevel) (Interpreter, error) {
	v, err := validator.New(mod)
	if err != nil {
		return nil, fmt.Errorf("New interpreter: \n\t%w", err)
	}
	if _, err := v.Validate(); err != nil {
		return nil, fmt.Errorf("New interpreter: \n\t%w", err)
	}
	inst, err := instance.New(mod)
	if err != nil {
		return nil, fmt.Errorf("New interpreter: \n\t%w", err)
	}
	stack := stack.New()
	return &interpreter{
		instance: inst,
		stack:    stack,
		cur:      &current{},
		debubber: debugger.New(debugLevel),
	}, nil
}

func (i *interpreter) Invoke(name string, locals []value.Value) ([]value.Value, error) {
	ext, err := i.instance.GetExport(name)
	if err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	if ext.ExternalValueType() != instance.ExternalValueTypeFunc {
		return nil, fmt.Errorf("Invoke: \n\t%w", FunctionIsRequired)
	}
	f := instance.GetExternVal[*instance.Function](ext)
	if err := validateLocals(f, locals); err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	if err := i.stack.Frame.Push(stack.Frame{Module: nil, Locals: nil}); err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	if err := i.stack.Label.Push(stack.Label{Instructions: nil, N: 0}); err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	for _, v := range locals {
		if err := i.stack.PushValue(v); err != nil {
			return nil, fmt.Errorf("Invoke: \n\t%w", err)
		}
	}
	if err := i.invokeFunction(f); err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	if err := i.execute(); err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	res, err := i.finishInvoke(f)
	if err != nil {
		return nil, fmt.Errorf("Invoke: \n\t%w", err)
	}
	return res, nil
}

// https://webassembly.github.io/spec/core/exec/instructions.html#invocation-of-function-address-a
func (i *interpreter) invokeFunction(f *instance.Function) error {
	// valudate local arguments and values on the stack
	if err := i.stack.Value.Validate(f.Type.Params); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	// get function arguments from the value stack
	locals, err := i.stack.PopValuesRev(len(f.Type.Params))
	if err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	if err := i.stack.Frame.Push(stack.Frame{Module: f.Module, Locals: locals}); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	if err := i.stack.Label.Push(stack.Label{Instructions: f.Code.Body, N: uint8(len(f.Type.Returns)), Sp: 0, Flag: true}); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	// sync current frame and label with top of the stack
	if err := i.cur.update(i.stack); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	// execute function instruction
	return nil
}

func (i *interpreter) execute() error {
	for {
		// frame := i.cur.frame
		label := i.cur.label
		contexSwitch := false
		if i.stack.Label.Len() <= 1 {
			return nil
		}
		for sp := label.Sp; sp < len(label.Instructions); sp++ {
			instr := label.Instructions[sp]
			res, err := i.step(instr)
			if err != nil {
				return fmt.Errorf("execute: %w", err)
			}
			label.Sp++
			switch res {
			case instructionResultTrap:
				return Trap
			case instructionResultCallFunc:
				if i.f == nil {
					return fmt.Errorf("execute: called function is not found")
				}
				if err := i.invokeFunction(i.f); err != nil {
					return fmt.Errorf("execute: %w", err)
				}
				contexSwitch = true
			case instructionResultReturn:
			case instructionResultEnterBlock:
				if err := i.cur.update(i.stack); err != nil {
					return fmt.Errorf("execute: %w", err)
				}
				contexSwitch = true
			case instructionResultLabelEnd:
				if i.isInvocationFinished() {
					return nil
				}
				contexSwitch = true
			case instructionResultRunNext:
				// go to next step
			}
			if contexSwitch {
				break
			}
		}
		// fmt.Println("function context is changed.")
	}
}

// https://webassembly.github.io/spec/core/exec/instructions.html#returning-from-a-function
func (i *interpreter) finishInvoke(f *instance.Function) ([]value.Value, error) {
	if err := i.stack.Value.Validate(f.Type.Returns); err != nil {
		return nil, fmt.Errorf("finish: %w", err)
	}
	values, err := i.stack.PopValuesRev(len(f.Type.Returns))
	if err != nil {
		return nil, fmt.Errorf("finish: %w", err)
	}
	i.debubber.ShowResult(values)
	return values, nil
}

func (i *interpreter) isInvocationFinished() bool {
	if i.stack.Label.Len() > 1 || i.stack.Label.IsEmpty() {
		return false
	}
	if i.stack.Frame.Len() == 1 {
		// frame stack: dummy
		fmt.Println("function invocation is finished.")
		return true
	}
	return false
}

func (i *interpreter) step(instr instruction.Instruction) (instructionResult, error) {
	switch instr.Opcode() {
	case instruction.NOP:
		return i.execNop(instr)
	case instruction.UNREACHABLE:
		return i.execUnreachable(instr)
	case instruction.DROP:
		return i.execDrop(instr)
	case instruction.SELECT:
		return i.execSelect(instr)
	case instruction.BLOCK:
		return i.execBlock(instr)
	case instruction.LOOP:
		return i.execLoop(instr)
	case instruction.IF:
		return i.execIf(instr)
	case instruction.BR:
		return i.execBr(instr)
	case instruction.BR_IF:
		return i.execBrIf(instr)
	case instruction.RETURN:
		return i.execReturn(instr)
	case instruction.I32_CONST, instruction.I64_CONST, instruction.F32_CONST, instruction.F64_CONST:
		return i.execConst(instr)
	case instruction.GET_LOCAL, instruction.SET_LOCAL, instruction.TEE_LOCAL:
		return i.execLocal(instr, i.cur.frame)
	case instruction.I32_ADD, instruction.I64_ADD, instruction.F32_ADD, instruction.F64_ADD,
		instruction.I32_SUB, instruction.I64_SUB,
		instruction.I32_MUL, instruction.I64_MUL,
		instruction.I32_DIV_S, instruction.I32_DIV_U, instruction.I64_DIV_S, instruction.I64_DIV_U,
		instruction.I32_REM_S, instruction.I32_REM_U, instruction.I64_REM_S, instruction.I64_REM_U,
		instruction.I32_AND, instruction.I64_AND,
		instruction.I32_OR, instruction.I64_OR,
		instruction.I32_XOR, instruction.I64_XOR,
		instruction.I32_SHL, instruction.I64_SHL:
		return i.execBinop(instr)
	case instruction.CALL:
		return i.execCall(instr)
	case instruction.END:
		return i.execLabelEnd(instr)
	default:
		// return instruction.InvalidOpcode
		return instructionResultTrap, nil
	}
}

func (i *interpreter) restoreStack() error {
	label, err := i.stack.Label.Pop()
	if err != nil {
		return fmt.Errorf("restore: %w", err)
	}
	if label.Flag {
		if _, err := i.stack.Frame.Pop(); err != nil {
			return fmt.Errorf("restore: %w", err)
		}
	}
	return i.cur.update(i.stack)
}

func (i *interpreter) expand(block types.BlockType) (*types.FuncType, error) {
	switch types.ValueType(block) {
	case types.I32, types.I64, types.F32, types.F64:
		return &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{types.ValueType(block)}}, nil
	case types.BLOCKTYPE:
		return &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{}}, nil
	default:
		if int(block) >= len(i.cur.frame.Module.Types) {
			return nil, fmt.Errorf("expand: function type is not found")
		}
		return i.cur.frame.Module.Types[int(block)], nil
	}
}

func validateLocals(f *instance.Function, locals []value.Value) error {
	params := f.Type.Params
	if len(params) != len(locals) {
		return fmt.Errorf("%w: expected=%d actual=%d", FunctionParamsDoesntMatch, len(params), len(locals))
	}
	for i, p := range params {
		l := locals[i]
		if l.ValType() != value.ValTypeNum || !p.IsNumber() {
			return FunctionParamTypesDoesntMatch
		}
	}
	return nil
}
