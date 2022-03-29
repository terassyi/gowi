package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

func TestBinop(t *testing.T) {
	for _, d := range []struct {
		interpreter *interpreter
		t           value.NumberType
		f           binopFunc
		exp         *stack.Stack
	}{
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.I32(1), value.I32(1)}, nil, nil)},
			f:           add,
			exp:         stackWithValueIgnoreError([]value.Value{value.I32(2)}, nil, nil),
		},
	} {
		err := d.interpreter.binop(d.t, d.f)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.interpreter.stack)
	}
}

func TestAdd_I32(t *testing.T) {
	for _, d := range []struct {
		a   value.I32
		b   value.I32
		exp value.I32
	}{
		{a: value.I32(0), b: value.I32(1), exp: value.I32(1)},
		{a: value.I32(0x0f), b: value.I32(0x1f), exp: value.I32(0x2e)},
	} {
		res, err := add(d.a, d.b)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}

func TestInterpreterLabelBlock(t *testing.T) {
	for _, d := range []struct {
		interpreter *interpreter
		ft          *types.FuncType
		exp         *stack.Label
	}{
		{
			interpreter: &interpreter{cur: &current{
				frame: nil,
				label: &stack.Label{
					Instructions: []instruction.Instruction{
						&instruction.Block{Imm: types.BlockType(0x40)},
						&instruction.Nop{},
						&instruction.End{},
					},
					N:  0,
					Sp: 0,
				},
			}},
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			interpreter: &interpreter{cur: &current{
				frame: nil,
				label: &stack.Label{
					Instructions: []instruction.Instruction{
						&instruction.Block{Imm: types.BlockType(types.I32)},
						&instruction.Nop{},
						&instruction.I32Const{Imm: 0},
						&instruction.End{},
					},
					N:  0,
					Sp: 0,
				},
			}},
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{types.I32}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.I32Const{Imm: 0},
					&instruction.End{},
				},
				N:  1,
				Sp: 0,
			},
		},
		{
			interpreter: &interpreter{cur: &current{
				frame: nil,
				label: &stack.Label{
					Instructions: []instruction.Instruction{
						&instruction.Block{Imm: types.BlockType(types.I32)},
						&instruction.Nop{},
						&instruction.End{},
						&instruction.Block{Imm: types.BlockType(types.I32)},
						&instruction.I32Const{Imm: 0},
						&instruction.End{},
						&instruction.End{},
					},
					N:  0,
					Sp: 3,
				},
			}},
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{types.I32}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.I32Const{Imm: 0},
					&instruction.End{},
				},
				N:  1,
				Sp: 0,
			},
		},
		{
			interpreter: &interpreter{cur: &current{
				frame: nil,
				label: &stack.Label{
					Instructions: []instruction.Instruction{
						&instruction.Block{Imm: types.BlockType(types.I32)},
						&instruction.Block{Imm: types.BlockType(types.I32)},
						&instruction.I32Const{Imm: 0},
						&instruction.End{},
						&instruction.End{},
						&instruction.End{},
					},
					N:  1,
					Sp: 0,
				},
			}},
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{types.I32}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Block{Imm: types.BlockType(types.I32)},
					&instruction.I32Const{Imm: 0},
					&instruction.End{},
					&instruction.End{},
				},
				N:  1,
				Sp: 0,
			},
		},
	} {
		label, _, err := d.interpreter.labelBlock(d.ft)
		require.NoError(t, err)
		assert.Equal(t, d.exp, label)
	}
}
