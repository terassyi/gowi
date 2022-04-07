package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"
)

func TestNewCode(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *code
	}{
		{
			payload: []byte{0x01, 0x07, 0x00, 0x20, 0x00, 0x20, 0x01, 0x6a, 0x0b},
			sec: &code{
				bodies: []*functionBody{
					{
						locals: []*localEntry{},
						code:   []instruction.Instruction{&instruction.GetLocal{Imm: uint32(0x00)}, &instruction.GetLocal{Imm: uint32(0x01)}, &instruction.I32Add{}, &instruction.End{}},
					},
				},
			},
		},
		// {
		// 	payload: []byte{0x03,
		// 		0x0f, 0x00, 0x41, 0x00, 0x41, 0x00, 0x2d, 0x00, 0x00, 0x41, 0x01, 0x6a, 0x3a, 0x00, 0x00, 0x0b,
		// 		0x08, 0x00, 0x41, 0x00, 0x2d, 0x00, 0x00, 0x0f, 0x0b,
		// 		0x08, 0x00, 0x10, 0x00, 0x10, 0x00, 0x10, 0x00, 0x0b,
		// 	},
		// 	sec: &Code{
		// 		Bodies: []*FunctionBody{
		// 			{
		// 				Locals: []*LocalEntry{},
		// 				Code:   []byte{0x41, 0x00, 0x41, 0x00, 0x2d, 0x00, 0x00, 0x41, 0x01, 0x6a, 0x3a, 0x00, 0x00, 0x0b},
		// 			},
		// 			{
		// 				Locals: []*LocalEntry{},
		// 				Code:   []byte{0x41, 0x00, 0x2d, 0x00, 0x00, 0x0f, 0x0b},
		// 			},
		// 			{
		// 				Locals: []*LocalEntry{},
		// 				// Code:   []byte{0x10, 0x00, 0x10, 0x00, 0x10, 0x00, 0x0b},
		// 				Code: []instruction.Instruction{&instruction.Call{Imm: uint32(0x00)}, &instruction.Call{Imm: uint32(0x00)}, &instruction.Call{Imm: uint32(0x00)}, &instruction.End{}},
		// 			},
		// 		},
		// 	},
		// },
		{
			payload: []byte{0x02,
				0x06, 0x01, 0x01, 0x7f, 0x20, 0x00, 0x0b,
				0x06, 0x01, 0x01, 0x7e, 0x20, 0x00, 0x0b,
			},
			sec: &code{
				bodies: []*functionBody{
					{
						locals: []*localEntry{{count: uint32(0x01), typ: types.I32}},
						// Code:   []byte{0x20, 0x00, 0x0b},
						code: []instruction.Instruction{&instruction.GetLocal{Imm: uint32(0x00)}, &instruction.End{}},
					},
					{
						locals: []*localEntry{{count: uint32(0x01), typ: types.I64}},
						// Code:   []byte{0x20, 0x00, 0x0b},
						code: []instruction.Instruction{&instruction.GetLocal{Imm: uint32(0x00)}, &instruction.End{}},
					},
				},
			},
		},
		{
			payload: []byte{0x01, 0x08, 0x00, 0x41, 0x00, 0x04, 0x40, 0x01, 0x0b, 0x0b},
			sec: &code{
				bodies: []*functionBody{
					{
						locals: []*localEntry{},
						code:   []instruction.Instruction{&instruction.I32Const{Imm: int32(0)}, &instruction.If{Imm: types.BlockType(types.ValueType(0x40))}, &instruction.Nop{}, &instruction.End{}, &instruction.End{}},
					},
				},
			},
		},
		{
			payload: []byte{0x01, 0x0d, 0x00, 0x41, 0x00, 0x04, 0x40, 0x41, 0x01, 0x04, 0x40, 0x01, 0x0b, 0x0b, 0x0b},
			sec: &code{
				bodies: []*functionBody{
					{
						locals: []*localEntry{},
						code: []instruction.Instruction{&instruction.I32Const{Imm: int32(0)}, &instruction.If{Imm: types.BlockType(types.ValueType(0x40))},
							&instruction.I32Const{Imm: int32(1)}, &instruction.If{Imm: types.BlockType(types.ValueType(0x40))},
							&instruction.Nop{}, &instruction.End{}, &instruction.End{}, &instruction.End{}},
					},
				},
			},
		},
	} {
		c, err := newCode(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, c)
	}
}
