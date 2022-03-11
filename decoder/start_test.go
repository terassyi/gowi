package decoder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStart(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *start
	}{
		{payload: []byte{0x02}, sec: &start{index: uint32(0x02)}},
	} {
		s, err := newStart(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, s)
	}
}
