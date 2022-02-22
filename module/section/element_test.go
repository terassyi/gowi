package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestElement(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Element
	}{
		{
			payload: []byte{0x01, 0x00, 0x41, 0x00, 0x0b, 0x01, 0x00},
			sec: &Element{
				count: 1,
				entries: []*ElementEntry{
					{
						Index:  0x00,
						Offset: []byte{0x41, 0x00},
						Number: uint32(0x01),
						Elems:  []uint32{0x00},
					},
				},
			},
		},
	} {
		e, err := NewElement(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec.count, e.count)
		assert.Equal(t, d.sec.entries, e.entries)
	}
}
