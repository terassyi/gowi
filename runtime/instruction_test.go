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

func TestBlockStackPush(t *testing.T) {
	for _, d := range []struct {
		bs  *blockStack
		op  instruction.Opcode
		exp *blockStack
	}{
		{
			bs:  newBlockStack(),
			op:  instruction.BLOCK,
			exp: &blockStack{inner: []block{&blockBlock{}}},
		},
		{
			bs:  newBlockStack(),
			op:  instruction.LOOP,
			exp: &blockStack{inner: []block{&blockLoop{}}},
		},
		{
			bs:  newBlockStack(),
			op:  instruction.IF,
			exp: &blockStack{inner: []block{&blockIf{}}},
		},
		{
			bs:  &blockStack{inner: []block{&blockIf{}}},
			op:  instruction.BLOCK,
			exp: &blockStack{inner: []block{&blockIf{}, &blockBlock{}}},
		},
	} {
		err := d.bs.push(d.op)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.bs)
	}
}

func TestBlockStackPop(t *testing.T) {
	for _, d := range []struct {
		bs   *blockStack
		exp  *blockStack
		expT blockType
	}{
		{
			bs:   &blockStack{inner: []block{&blockBlock{}}},
			exp:  &blockStack{inner: []block{}},
			expT: blockTypeBlock,
		},
		{
			bs:   &blockStack{inner: []block{&blockIf{}, &blockBlock{}}},
			exp:  &blockStack{inner: []block{&blockIf{}}},
			expT: blockTypeBlock,
		},
	} {
		v, err := d.bs.pop()
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.bs)
		assert.Equal(t, d.expT, v.typ())
	}
}
func TestInterpreterLabelBlock(t *testing.T) {
	for _, d := range []struct {
		interpreter *interpreter
		op          instruction.Opcode
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
			op: instruction.BLOCK,
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
			op: instruction.BLOCK,
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
			op: instruction.BLOCK,
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
			op: instruction.BLOCK,
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
		{
			interpreter: &interpreter{cur: &current{
				frame: nil,
				label: &stack.Label{
					Instructions: []instruction.Instruction{
						&instruction.If{Imm: types.BlockType(0x40)},
						&instruction.Nop{},
						&instruction.End{},
					},
					N:  0,
					Sp: 0,
				},
			}},
			op: instruction.BLOCK,
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
						&instruction.If{Imm: types.BlockType(0x40)},
						&instruction.Nop{},
						&instruction.Else{},
						&instruction.End{},
					},
					N:  0,
					Sp: 0,
				},
			}},
			op: instruction.BLOCK,
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.Else{},
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
						&instruction.If{Imm: types.BlockType(0x40)},
						&instruction.Nop{},
						&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)},
						&instruction.End{},
						&instruction.Else{},
						&instruction.End{},
					},
					N:  0,
					Sp: 0,
				},
			}},
			op: instruction.BLOCK,
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.End{},
					&instruction.Else{},
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
						&instruction.If{Imm: types.BlockType(0x40)},
						&instruction.Nop{},
						&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)},
						&instruction.End{},
						&instruction.Else{},
						&instruction.Block{Imm: types.BlockType(types.BLOCKTYPE)},
						&instruction.Nop{},
						&instruction.End{},
						&instruction.End{},
						&instruction.End{},
					},
					N:  0,
					Sp: 0,
				},
			}},
			op: instruction.BLOCK,
			ft: &types.FuncType{Params: types.ResultType{}, Returns: types.ResultType{}},
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.End{},
					&instruction.Else{},
					&instruction.Block{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.Nop{},
					&instruction.End{},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
	} {
		label, _, err := d.interpreter.labelBlock(d.ft)
		require.NoError(t, err)
		assert.Equal(t, d.exp, label)
	}
}

func TestIfElseLabelBlock(t *testing.T) {
	for _, d := range []struct {
		label *stack.Label
		cond  value.I32
		exp   *stack.Label
	}{
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(1),
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{},
				N:            0,
				Sp:           0,
			},
			cond: value.I32(0),
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(0),
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Nop{},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(1),
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
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.I32Const{Imm: 7},
					&instruction.Else{},
					&instruction.I32Const{Imm: 7},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(1),
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.I32Const{Imm: 7},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.I32Const{Imm: 7},
					&instruction.Else{},
					&instruction.I32Const{Imm: 8},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(0),
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.I32Const{Imm: 8},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.GetLocal{Imm: 1},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.Call{Imm: 0},
					&instruction.Block{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.End{},
					&instruction.Nop{},
					&instruction.End{},
					&instruction.Else{},
					&instruction.Call{Imm: 0},
					&instruction.Block{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.End{},
					&instruction.Nop{},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(1),
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.GetLocal{Imm: 1},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.Call{Imm: 0},
					&instruction.Block{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.End{},
					&instruction.Nop{},
					&instruction.End{},
					&instruction.End{},
				},
				N:  0,
				Sp: 0,
			},
		},
		{
			label: &stack.Label{
				Instructions: []instruction.Instruction{
					// if0
					&instruction.GetLocal{Imm: 1},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)}, // if1
					&instruction.Call{Imm: 0},
					&instruction.Block{Imm: types.BlockType(types.BLOCKTYPE)},
					&instruction.End{}, // end of block
					&instruction.Nop{},
					&instruction.End{},  // end of if1
					&instruction.Else{}, // else of if0
					&instruction.Call{Imm: 0},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)}, // end of if2
					&instruction.Nop{},
					&instruction.Else{}, // else of if2
					&instruction.Call{Imm: 0},
					&instruction.Call{Imm: 0},
					&instruction.End{}, // end of if2
					&instruction.End{}, // end of if0
				},
				N:  0,
				Sp: 0,
			},
			cond: value.I32(0), //e lse
			exp: &stack.Label{
				Instructions: []instruction.Instruction{
					&instruction.Call{Imm: 0},
					&instruction.If{Imm: types.BlockType(types.BLOCKTYPE)}, // end of if2
					&instruction.Nop{},
					&instruction.Else{}, // else of if2
					&instruction.Call{Imm: 0},
					&instruction.Call{Imm: 0},
					&instruction.End{}, // end of if2
					&instruction.End{}, // end of if0
				},
				N:  0,
				Sp: 0,
			},
		},
	} {
		condLabel, err := ifElseLabel(d.label, d.cond)
		require.NoError(t, err)
		assert.Equal(t, d.exp, condLabel)
	}
}
