package section

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewType(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sec     *Type
	}{
		{payload: []byte{0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f}, sec: &Type{count: 1}},
		{payload: []byte{0x01, 0x60, 0x00, 0x01, 0x7f, 0x01}, sec: &Type{count: 1}},
		{payload: []byte{0x01, 0x60, 0x00, 0x00}, sec: &Type{count: 1}},
		{payload: []byte{0x02, 0x60, 0x01, 0x7f, 0x00, 0x60, 0x00, 0x00}, sec: &Type{count: 2}},
		{payload: []byte{0x02, 0x60, 0x00, 0x01, 0x7f, 0x60, 0x01, 0x7f, 0x01, 0x7f}, sec: &Type{count: 2}},
	} {
		typ, err := NewType(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.sec.count, typ.count)
	}
}
