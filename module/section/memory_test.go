package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestMemory(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Memory
	}{
		{
			payload: []byte{0x01, 0x00, 0x01},
			sec: &Memory{
				Entries: []*types.MemoryType{
					{
						Limits: &types.ResizableLimits{
							Flag:    false,
							Initial: uint32(0x01),
						},
					},
				},
			},
		},
	} {
		m, err := NewMemory(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, m)
	}
}
