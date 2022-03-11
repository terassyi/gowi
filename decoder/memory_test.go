package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestMemory(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *memory
	}{
		{
			payload: []byte{0x01, 0x00, 0x01},
			sec: &memory{
				entries: []*types.MemoryType{
					{
						Limits: &types.Limits{
							Min: uint32(0x01),
							Max: uint32(0x00),
						},
					},
				},
			},
		},
	} {
		m, err := newMemory(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, m)
	}
}
