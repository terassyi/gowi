package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeFuncType(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		f       *FuncType
		n       int
	}{
		{payload: []byte{0x02, 0x7f, 0x7f, 0x01, 0x7f}, f: &FuncType{Params: []ValueType{I32, I32}, Returns: []ValueType{I32}}, n: 5},
		{payload: []byte{0x00, 0x01, 0x7f}, f: &FuncType{Params: []ValueType{}, Returns: []ValueType{I32}}, n: 3},
		{payload: []byte{0x00, 0x00}, f: &FuncType{Params: []ValueType{}, Returns: []ValueType{}}, n: 2},
		{payload: []byte{0x01, 0x60, 0x00}, f: &FuncType{Params: []ValueType{FUNC}, Returns: []ValueType{}}, n: 3},
	} {
		f, n, err := DecodeFuncType(d.payload)
		require.NoError(t, err)
		assert.Equal(t, d.n, n)
		assert.Equal(t, d.f, f)
	}
}
