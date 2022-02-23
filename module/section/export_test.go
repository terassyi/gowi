package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/types"
)

func TestNewExport(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Export
	}{
		{
			payload: []byte{0x01, 0x03, 0x61, 0x64, 0x64, 0x00, 0x00},
			sec: &Export{
				Entries: []*ExportEntry{
					{
						FieldString: []byte{0x61, 0x64, 0x64},
						Kind:        types.EXTERNAL_KIND_FUNCTION,
						Index:       0x00,
					},
				},
			},
		},
		{
			payload: []byte{0x01, 0x0e, 0x67, 0x65, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x50, 0x6c, 0x75, 0x73, 0x31, 0x00, 0x01},
			sec: &Export{
				Entries: []*ExportEntry{
					{
						FieldString: []byte{0x67, 0x65, 0x74, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x50, 0x6c, 0x75, 0x73, 0x31},
						Kind:        types.EXTERNAL_KIND_FUNCTION,
						Index:       uint32(0x01),
					},
				},
			},
		},
	} {
		e, err := NewExport(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, e)
	}
}
