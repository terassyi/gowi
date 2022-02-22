package section

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSectionNew(t *testing.T) {
	for _, d := range []struct {
		id      uint8
		payload []byte
		sec     *Type
	}{
		{id: 0x01, payload: []byte{0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f}, sec: &Type{count: 1}},
	} {
		sec, err := New(d.id, d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec.count, sec.(*Type).count)
	}
}
func TestSectionNew_InvalidSectionCode(t *testing.T) {
	_, err := New(0xf, []byte{})
	if !errors.Is(err, InvalidSectionCode) {
		t.Error(err)
	}
}
