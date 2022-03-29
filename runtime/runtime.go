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
	"github.com/terassyi/gowi/validator"
)

var (
	FunctionIsRequired            error = errors.New("External value type function is required")
	FunctionParamsDoesntMatch     error = errors.New("Number of function parameters doesn't match")
	FunctionParamTypesDoesntMatch error = errors.New("Function parameter type doesn't match")
)

type Interpreter interface {
	Invoke(string, []value.Value) error
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
	read  int
}

func (c *current) update(s *stack.Stack) error {
	frame, err := s.Frame.Top()
	if err != nil {
		return err
	}
	c.frame = frame
	label, err := s.Label.Top()
	if err != nil {
		return err
	}
	c.label = label
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

func (i *interpreter) Invoke(name string, locals []value.Value) error {
	ext, err := i.instance.GetExport(name)
	if err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	if ext.ExternalValueType() != instance.ExternalValueTypeFunc {
		return fmt.Errorf("Invoke: \n\t%w", FunctionIsRequired)
	}
	f := instance.GetExternVal[*instance.Function](ext)
	if err := validateLocals(f, locals); err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	if err := i.stack.Frame.Push(stack.Frame{Module: nil, Locals: nil}); err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	if err := i.stack.Label.Push(stack.Label{Instructions: nil, N: 0}); err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	for _, v := range locals {
		if err := i.stack.Value.Push(v); err != nil {
			return fmt.Errorf("Invoke: \n\t%w", err)
		}
	}
	if err := i.invokeFunction(f); err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	if err := i.finishInvoke(f); err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	return nil
}

// https://webassembly.github.io/spec/core/exec/instructions.html#invocation-of-function-address-a
func (i *interpreter) invokeFunction(f *instance.Function) error {
	// valudate local arguments and values on the stack
	if err := i.stack.Value.Validate(f.Type.Params); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	// get function arguments from the value stack
	locals, err := i.stack.Value.PopNRev(len(f.Type.Params))
	if err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	if err := i.stack.Frame.Push(stack.Frame{Module: f.Module, Locals: locals}); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	if err := i.stack.Label.Push(stack.Label{Instructions: f.Code.Body, N: uint8(len(f.Type.Returns))}); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	// sync current frame and label with top of the stack
	if err := i.cur.update(i.stack); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	// execute function instruction
	if err := i.execute(); err != nil {
		return fmt.Errorf("Invoke function: %w", err)
	}
	return nil
}

func (i *interpreter) execute() error {
	for _, instr := range i.cur.label.Instructions {
		res, err := i.step(instr)
		if err != nil {
			return fmt.Errorf("execute: %w", err)
		}
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
		case instructionResultReturn:
		case instructionResultEnterBlock:
		case instructionResultLabelEnd:
			if i.isInvocationFinished() {
				return nil
			}
		case instructionResultRunNext:
			// go to next step
		}
	}
	return nil
}

// https://webassembly.github.io/spec/core/exec/instructions.html#returning-from-a-function
func (i *interpreter) finishInvoke(f *instance.Function) error {
	if err := i.stack.Value.Validate(f.Type.Returns); err != nil {
		return fmt.Errorf("finish: %w", err)
	}
	values, err := i.stack.Value.PopNRev(len(f.Type.Returns))
	if err != nil {
		return fmt.Errorf("finish: %w", err)
	}
	i.debubber.ShowResult(values)
	return nil
}

func (i *interpreter) isInvocationFinished() bool {
	if i.stack.Label.Len() > 1 || i.stack.Label.IsEmpty() {
		return false
	}
	if i.stack.Frame.Len() == 1 {
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
	case instruction.I32_CONST, instruction.I64_CONST, instruction.F32_CONST, instruction.F64_CONST:
		return i.execConst(instr)
	case instruction.GET_LOCAL, instruction.SET_LOCAL, instruction.TEE_LOCAL:
		return i.execLocal(instr, i.cur.frame)
	case instruction.I32_ADD, instruction.I64_ADD, instruction.F32_ADD, instruction.F64_ADD:
		return i.execBinop(instr)
	case instruction.CALL:
		return i.execCall(instr)
	case instruction.END:
		return i.labelEnd(instr)
	default:
		// return instruction.InvalidOpcode
		return instructionResultTrap, nil
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
