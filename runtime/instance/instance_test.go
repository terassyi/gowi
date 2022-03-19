package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/runtime/value"
)

func TestGetVal_I32(t *testing.T) {
	for _, d := range []struct {
		val      value.Value
		expected value.I32
	}{
		{val: value.I32(0), expected: value.I32(0)},
		{val: value.I32(1), expected: value.I32(1)},
		{val: value.I32(-1), expected: value.I32(-1)},
	} {
		v := GetVal[value.I32](d.val)
		assert.Equal(t, d.expected, v)
	}
}

func TestGetVal_F64(t *testing.T) {
	for _, d := range []struct {
		val      value.Value
		expected value.F64
	}{
		{val: value.F64(0), expected: value.F64(0)},
		{val: value.F64(1.1), expected: value.F64(1.1)},
		{val: value.F64(-1.0012), expected: value.F64(-1.0012)},
	} {
		v := GetVal[value.F64](d.val)
		assert.Equal(t, d.expected, v)
	}
}

func TestGetVal_Panic(t *testing.T) {
	require.Panics(t, func() {
		GetVal[value.I64](value.I32(0))
	})
}
