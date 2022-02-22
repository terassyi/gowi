package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStart(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Start
	}{
		{payload: []byte{0x02}, sec: &Start{Index: uint32(0x02)}},
	} {
		s, err := NewStart(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec, s)
	}
}
