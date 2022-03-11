package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *data
	}{
		{
			payload: []byte{0x01, 0x00, 0x41, 0x00, 0x0b, 0x01, 0x41},
			sec: &data{
				entries: []*dataSegment{
					{
						index:  uint32(0x00),
						offset: []byte{0x41, 0x00},
						size:   uint32(0x01),
						data:   []byte{0x41},
					},
				},
			},
		},
		{
			payload: []byte{0x05,
				0x00, 0x41, 0x00, 0x0b, 0x01, 0x61,
				0x00, 0x41, 0x03, 0x0b, 0x01, 0x62,
				0x00, 0x41, 0xe4, 0x00, 0x0b, 0x03, 0x63, 0x64, 0x65,
				0x00, 0x41, 0x05, 0x0b, 0x01, 0x78,
				0x00, 0x41, 0x03, 0x0b, 0x01, 0x63},
			sec: &data{
				entries: []*dataSegment{
					{
						index:  uint32(0x00),
						offset: []byte{0x41, 0x00},
						size:   uint32(0x01),
						data:   []byte{0x61},
					},
					{
						index:  uint32(0x00),
						offset: []byte{0x41, 0x03},
						size:   uint32(0x01),
						data:   []byte{0x62},
					},
					{
						index:  uint32(0x00),
						offset: []byte{0x41, 0xe4, 0x00},
						size:   uint32(0x03),
						data:   []byte{0x63, 0x64, 0x65},
					},
					{
						index:  uint32(0x00),
						offset: []byte{0x41, 0x05},
						size:   uint32(0x01),
						data:   []byte{0x78},
					},
					{
						index:  uint32(0x00),
						offset: []byte{0x41, 0x03},
						size:   uint32(0x01),
						data:   []byte{0x63},
					},
				},
			},
		},
	} {
		data, err := newData(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, data)
	}
}
