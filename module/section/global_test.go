package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewGloabal(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Global
	}{
		{
			payload: []byte{0x02, 0x7f, 0x00, 0x41, 0x7e, 0x0b, 0x7f, 0x01, 0x41, 0x74, 0x0b},
			sec: &Global{
				Globals: []*GlobalEntry{
					{
						Type: &types.GlobalType{
							ContentType: types.I32,
							Mut:         false,
						},
						Init: []byte{0x41, 0x7e},
					},
					{
						Type: &types.GlobalType{
							ContentType: types.I32,
							Mut:         true,
						},
						Init: []byte{0x41, 0x74},
					},
				},
			},
		},
	} {
		g, err := NewGlobal(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, g)
	}
}
