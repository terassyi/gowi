package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewImport(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *imports
	}{
		{
			payload: []byte{0x01, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x03, 0x6c, 0x6f, 0x67, 0x00, 0x00},
			sec: &imports{
				entries: []*importEntry{
					{
						moduleNameLength: uint32(0x07),
						moduleName:       []byte{0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65},
						fieldLength:      uint32(0x03),
						fieldString:      []byte{0x6c, 0x6f, 0x67},
						kind:             types.EXTERNAL_KIND_FUNCTION,
						typ:              types.VarUint32(0x00),
					},
				},
			},
		},
		{
			payload: []byte{0x02,
				0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65,
				0x03, 0x6c, 0x6f, 0x67,
				0x00, 0x00,
				0x02, 0x6a, 0x73,
				0x03, 0x6d, 0x65, 0x6d,
				0x02, 0x00, 0x01},
			sec: &imports{
				entries: []*importEntry{
					{
						moduleNameLength: uint32(0x07),
						moduleName:       []byte{0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65},
						fieldLength:      uint32(0x03),
						fieldString:      []byte{0x6c, 0x6f, 0x67},
						kind:             types.EXTERNAL_KIND_FUNCTION,
						typ:              types.VarUint32(0x00),
					},
					{
						moduleNameLength: uint32(0x02),
						moduleName:       []byte{0x6a, 0x73},
						fieldLength:      uint32(0x03),
						fieldString:      []byte{0x6d, 0x65, 0x6d},
						kind:             types.EXTERNAL_KIND_MEMORY,
						typ: &types.MemoryType{
							Limits: &types.Limits{
								Min: uint32(0x01),
								Max: uint32(0x00),
							},
						},
					},
				},
			},
		},
		{
			payload: []byte{0x02,
				0x02, 0x6a, 0x73,
				0x06, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79,
				0x02, 0x00, 0x01,
				0x02, 0x6a, 0x73,
				0x05, 0x74, 0x61, 0x62, 0x6c, 0x65,
				0x01, 0x70, 0x00, 0x01},
			sec: &imports{
				entries: []*importEntry{
					{
						moduleNameLength: uint32(0x02),
						moduleName:       []byte{0x6a, 0x73},
						fieldLength:      uint32(0x06),
						fieldString:      []byte{0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79},
						kind:             types.EXTERNAL_KIND_MEMORY,
						typ: &types.MemoryType{
							Limits: &types.Limits{
								Min: uint32(0x01),
								Max: uint32(0x00),
							},
						},
					},
					{
						moduleNameLength: uint32(0x02),
						moduleName:       []byte{0x6a, 0x73},
						fieldLength:      uint32(0x05),
						fieldString:      []byte{0x74, 0x61, 0x62, 0x6c, 0x65},
						kind:             types.EXTERNAL_KIND_TABLE,
						typ: &types.TableType{
							ElementType: types.ElemType(types.ANYFUNC),
							Limits: &types.Limits{
								Min: uint32(0x01),
								Max: uint32(0x00),
							},
						},
					},
				},
			},
		},
	} {
		i, err := newImport(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, i)
	}
}
