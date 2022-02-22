package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFunction(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Function
	}{
		{payload: []byte{0x02, 0x01, 0x00}, sec: &Function{count: 2, types: []uint32{0x01, 0x00}}},
		{payload: []byte{0x02, 0x00, 0x00}, sec: &Function{count: 2, types: []uint32{0x00, 0x00}}},
	} {
		f, err := NewFunction(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, f)
	}
}
