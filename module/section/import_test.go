package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewImport(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Import
	}{
		{
			payload: []byte{0x01, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x03, 0x6c, 0x6f, 0x67, 0x00, 0x00},
			sec: &Import{
				count: 1,
				entries: []*ImportEntry{
					{
						ModuleNameLength: uint32(0x07),
						ModuleName:       []byte{0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65},
						FieldLength:      uint32(0x03),
						FieldString:      []byte{0x6c, 0x6f, 0x67},
						Kind:             types.EXTERNAL_KIND_FUNCTION,
						Type:             types.VarUint32(0x00),
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
			sec: &Import{
				count: 2,
				entries: []*ImportEntry{
					{
						ModuleNameLength: uint32(0x07),
						ModuleName:       []byte{0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65},
						FieldLength:      uint32(0x03),
						FieldString:      []byte{0x6c, 0x6f, 0x67},
						Kind:             types.EXTERNAL_KIND_FUNCTION,
						Type:             types.VarUint32(0x00),
					},
					{
						ModuleNameLength: uint32(0x02),
						ModuleName:       []byte{0x6a, 0x73},
						FieldLength:      uint32(0x03),
						FieldString:      []byte{0x6d, 0x65, 0x6d},
						Kind:             types.EXTERNAL_KIND_MEMORY,
						Type: &types.MemoryType{
							Limits: &types.ResizableLimits{
								Flag:    false,
								Initial: uint32(0x01),
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
			sec: &Import{
				count: 2,
				entries: []*ImportEntry{
					{
						ModuleNameLength: uint32(0x02),
						ModuleName:       []byte{0x6a, 0x73},
						FieldLength:      uint32(0x06),
						FieldString:      []byte{0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79},
						Kind:             types.EXTERNAL_KIND_MEMORY,
						Type: &types.MemoryType{
							Limits: &types.ResizableLimits{
								Flag:    false,
								Initial: uint32(0x01),
							},
						},
					},
					{
						ModuleNameLength: uint32(0x02),
						ModuleName:       []byte{0x6a, 0x73},
						FieldLength:      uint32(0x05),
						FieldString:      []byte{0x74, 0x61, 0x62, 0x6c, 0x65},
						Kind:             types.EXTERNAL_KIND_TABLE,
						Type: &types.TableType{
							ElementType: types.ElemType(types.ANYFUNC),
							Limits: &types.ResizableLimits{
								Flag:    false,
								Initial: uint32(0x01),
							},
						},
					},
				},
			},
		},
	} {
		i, err := NewImport(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec.count, i.count)
		assert.Equal(t, d.sec.entries, i.entries)
	}
}
