package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewTable(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *table
	}{
		{
			payload: []byte{0x01, 0x70, 0x00, 0x02},
			sec: &table{
				entries: []*types.TableType{
					{
						ElementType: types.ElemTypeFuncref,
						Limits:      &types.Limits{Min: uint32(0x02), Max: uint32(0x00)},
					},
				},
			},
		},
	} {
		ty, err := newTable(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, ty)
	}
}
