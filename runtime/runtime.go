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
}

type interpreter struct {
	instance *instance.Module
	stack    *stack.Stack
	debubber *debugger.Debugger
	f        *instance.Function
}

type functionContext struct {
}

// instanciate an interpreter
// https://webassembly.github.io/spec/core/exec/modules.html#instantiation
func New(mod *structure.Module, externalvals []instance.ExternalValue) (Interpreter, error) {
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
	i.f = f
	if err := i.stack.Frame.Push(stack.Frame{Module: nil, Locals: nil}); err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	for _, v := range locals {
		if err := i.stack.Value.Push(v); err != nil {
			return fmt.Errorf("Invoke: \n\t%w", err)
		}
	}
	res, err := i.call()
	if err != nil {
		return fmt.Errorf("Invoke: \n\t%w", err)
	}
	fmt.Printf("Invocation Result %s() = %v\n\n", name, res)
	return nil
}

func (i *interpreter) call() (value.Result, error) {
	for _, instr := range i.f.Code.Body {
		if err := i.step(instr); err != nil {
			return nil, fmt.Errorf("Invocation call: %w", err)
		}
	}
	return nil, nil
}

func (i *interpreter) step(instr instruction.Instruction) error {

	return nil
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
