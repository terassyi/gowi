package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestData(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Data
	}{
		{
			payload: []byte{0x01, 0x00, 0x41, 0x00, 0x0b, 0x01, 0x41},
			sec: &Data{
				count: uint32(0x01),
				entries: []*DataSegment{
					{
						Index:  uint32(0x00),
						Offset: []byte{0x41, 0x00},
						Size:   uint32(0x01),
						Data:   []byte{0x41},
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
			sec: &Data{
				count: uint32(0x05),
				entries: []*DataSegment{
					{
						Index:  uint32(0x00),
						Offset: []byte{0x41, 0x00},
						Size:   uint32(0x01),
						Data:   []byte{0x61},
					},
					{
						Index:  uint32(0x00),
						Offset: []byte{0x41, 0x03},
						Size:   uint32(0x01),
						Data:   []byte{0x62},
					},
					{
						Index:  uint32(0x00),
						Offset: []byte{0x41, 0xe4, 0x00},
						Size:   uint32(0x03),
						Data:   []byte{0x63, 0x64, 0x65},
					},
					{
						Index:  uint32(0x00),
						Offset: []byte{0x41, 0x05},
						Size:   uint32(0x01),
						Data:   []byte{0x78},
					},
					{
						Index:  uint32(0x00),
						Offset: []byte{0x41, 0x03},
						Size:   uint32(0x01),
						Data:   []byte{0x63},
					},
				},
			},
		},
	} {
		data, err := NewData(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec.count, data.count)
		assert.Equal(t, d.sec.entries, data.entries)
	}
}
