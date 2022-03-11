package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewGloabal(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *global
	}{
		{
			payload: []byte{0x02, 0x7f, 0x00, 0x41, 0x7e, 0x0b, 0x7f, 0x01, 0x41, 0x74, 0x0b},
			sec: &global{
				globals: []*globalEntry{
					{
						typ: &types.GlobalType{
							ContentType: types.I32,
							Mut:         false,
						},
						init: []byte{0x41, 0x7e},
					},
					{
						typ: &types.GlobalType{
							ContentType: types.I32,
							Mut:         true,
						},
						init: []byte{0x41, 0x74},
					},
				},
			},
		},
	} {
		g, err := newGlobal(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, g)
	}
}
