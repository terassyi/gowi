package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
)

func TestBinop(t *testing.T) {
	for _, d := range []struct {
		interpreter *interpreter
		t           value.NumberType
		f           binopFunc
		exp         *stack.Stack
	}{
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.I32(1), value.I32(1)}, nil, nil)},
			f:           add,
			exp:         stackWithValueIgnoreError([]value.Value{value.I32(2)}, nil, nil),
		},
	} {
		err := d.interpreter.binop(d.t, d.f)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.interpreter.stack)
	}
}

func TestAdd_I32(t *testing.T) {
	for _, d := range []struct {
		a   value.I32
		b   value.I32
		exp value.I32
	}{
		{a: value.I32(0), b: value.I32(1), exp: value.I32(1)},
		{a: value.I32(0x0f), b: value.I32(0x1f), exp: value.I32(0x2e)},
	} {
		res, err := add(d.a, d.b)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}
