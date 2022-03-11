package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElement(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *element
	}{
		{
			payload: []byte{0x01, 0x00, 0x41, 0x00, 0x0b, 0x01, 0x00},
			sec: &element{
				entries: []*elementEntry{
					{
						index:  0x00,
						offset: []byte{0x41, 0x00},
						number: uint32(0x01),
						elems:  []uint32{0x00},
					},
				},
			},
		},
	} {
		e, err := newElement(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, e)
	}
}
