package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewTable(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Table
	}{
		{
			payload: []byte{0x01, 0x70, 0x00, 0x02},
			sec: &Table{
				Entries: []*types.TableType{
					{
						ElementType: types.ElemType(types.ANYFUNC),
						Limits:      &types.ResizableLimits{Flag: false, Initial: uint32(0x02)},
					},
				},
			},
		},
	} {
		ty, err := NewTable(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, ty)
	}
}
