package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewType(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *typ
	}{
		{
			payload: []byte{0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f},
			sec: &typ{
				entries: []*types.FuncType{
					{
						Params:  []types.ValueType{types.I32, types.I32},
						Returns: []types.ValueType{types.I32},
					},
				},
			},
		},
		{
			payload: []byte{0x01, 0x60, 0x00, 0x01, 0x7f},
			sec: &typ{
				entries: []*types.FuncType{
					{
						Params:  []types.ValueType{},
						Returns: []types.ValueType{types.I32},
					},
				},
			},
		},
		{
			payload: []byte{0x01, 0x60, 0x00, 0x00},
			sec: &typ{
				entries: []*types.FuncType{
					{
						Params:  []types.ValueType{},
						Returns: []types.ValueType{},
					},
				},
			},
		},
		{
			payload: []byte{0x02, 0x60, 0x01, 0x7f, 0x00, 0x60, 0x00, 0x00},
			sec: &typ{
				entries: []*types.FuncType{
					{
						Params:  []types.ValueType{types.I32},
						Returns: []types.ValueType{},
					},
					{
						Params:  []types.ValueType{},
						Returns: []types.ValueType{},
					},
				},
			},
		},
		{
			payload: []byte{0x02, 0x60, 0x00, 0x01, 0x7f, 0x60, 0x01, 0x7f, 0x01, 0x7f},
			sec: &typ{
				entries: []*types.FuncType{
					{
						Params:  []types.ValueType{},
						Returns: []types.ValueType{types.I32},
					},
					{
						Params:  []types.ValueType{types.I32},
						Returns: []types.ValueType{types.I32},
					},
				},
			},
		},
	} {
		typ, err := newType(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, typ)
	}
}
