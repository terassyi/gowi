package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/decoder"
	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/debugger"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
	"github.com/terassyi/gowi/validator"
)

func TestValidateLocals(t *testing.T) {
	for _, d := range []struct {
		f      *instance.Function
		locals []value.Value
	}{
		{f: &instance.Function{Type: &types.FuncType{Params: []types.ValueType{types.I32, types.F32}, Returns: []types.ValueType{types.F64}}}, locals: []value.Value{value.I32(0), value.F32(1.0)}},
		{f: &instance.Function{Type: &types.FuncType{Params: []types.ValueType{types.I32}, Returns: []types.ValueType{types.F64}}}, locals: []value.Value{value.I32(0)}},
		{f: &instance.Function{Type: &types.FuncType{Params: []types.ValueType{}, Returns: []types.ValueType{}}}, locals: []value.Value{}},
	} {
		err := validateLocals(d.f, d.locals)
		require.NoError(t, err)
	}
}

func TestInvoke(t *testing.T) {
	for _, d := range []struct {
		path   string
		export string
		args   []value.Value
		exp    []value.Value
	}{
		{path: "../examples/func1.wasm", export: "add", args: []value.Value{value.I32(0), value.I32(2)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/func1.wasm", export: "add", args: []value.Value{value.I32(13), value.I32(27)}, exp: []value.Value{value.I32(40)}},
		{path: "../examples/call_func1.wasm", export: "getAnswerPlus1", args: []value.Value{}, exp: []value.Value{value.I32(43)}},
		{path: "../examples/call_func_nested.wasm", export: "rootFunc", args: []value.Value{value.I32(1), value.I32(2)}, exp: []value.Value{value.I32(6)}},
		{path: "../examples/block.wasm", export: "singular", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/block.wasm", export: "multi", args: []value.Value{}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/block.wasm", export: "nest", args: []value.Value{}, exp: []value.Value{value.I32(150)}},
		{path: "../examples/block.wasm", export: "deep", args: []value.Value{}, exp: []value.Value{value.I32(150)}},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := instance.New(mod)
		require.NoError(t, err)
		interpreter := &interpreter{
			instance: ins,
			stack:    stack.New(),
			cur:      &current{},
			debubber: debugger.New(debugger.DebugLevelLogOnlyStdout),
		}
		res, err := interpreter.Invoke(d.export, d.args)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}

func TestStep(t *testing.T) {
	for _, d := range []struct {
		interpreter *interpreter
		instr       instruction.Instruction
		exp         any
		expStack    *stack.Stack
		expCur      *current
	}{
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}), cur: nil},
			instr:       &instruction.I32Const{Imm: 0},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I32(0)}, []stack.Frame{}, []stack.Label{}),
			expCur:      nil,
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.GetLocal{Imm: 0},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I32(0)}, []stack.Frame{}, []stack.Label{}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.GetLocal{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.F64(0.1)}, []stack.Frame{}, []stack.Label{}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.F64(1.5)}, []stack.Frame{}, []stack.Label{}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.SetLocal{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(1.5)}}, label: nil},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.F64(9.0)}, []stack.Frame{}, []stack.Label{}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.TeeLocal{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.F64(9.0)}, []stack.Frame{}, []stack.Label{}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(9.0)}}, label: nil},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.I32(0xf0)}, []stack.Frame{}, []stack.Label{}), cur: nil},
			instr:       &instruction.Drop{},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}),
			expCur:      nil,
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.I32(0xff), value.I32(0xee), value.I32(0x0)}, []stack.Frame{}, []stack.Label{}), cur: nil},
			instr:       &instruction.Select{},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I32(0xee)}, []stack.Frame{}, []stack.Label{}),
			expCur:      nil,
		},
	} {
		_, err := d.interpreter.step(d.instr)
		require.NoError(t, err)
		assert.Equal(t, d.expStack, d.interpreter.stack)
		assert.Equal(t, d.expCur, d.interpreter.cur)
	}
}

func stackWithValueIgnoreError(values []value.Value, frames []stack.Frame, labels []stack.Label) *stack.Stack {
	s, _ := stack.WithValue(values, frames, labels)
	return s
}
