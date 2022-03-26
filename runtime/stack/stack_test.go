package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/runtime/value"
)

func TestValueStackPush(t *testing.T) {
	for _, d := range []struct {
		s   *ValueStack
		val value.Value
		exp *ValueStack
	}{
		{s: &ValueStack{values: []value.Value{}, sp: 0}, val: value.I32(0), exp: &ValueStack{values: []value.Value{value.I32(0)}, sp: 0}},
		{s: &ValueStack{values: []value.Value{value.I32(0)}, sp: 0}, val: value.I32(0xff), exp: &ValueStack{values: []value.Value{value.I32(0), value.I32(0xff)}, sp: 0}},
		{s: &ValueStack{values: []value.Value{value.I32(0)}, sp: 0}, val: value.F32(1.1), exp: &ValueStack{values: []value.Value{value.I32(0), value.F32(1.1)}, sp: 0}},
	} {
		err := d.s.Push(d.val)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.s)
	}
}

func TestValueStackPop(t *testing.T) {
	for _, d := range []struct {
		s      *ValueStack
		expVal value.Value
		expS   *ValueStack
	}{
		{s: &ValueStack{values: []value.Value{value.I32(0)}, sp: 0}, expVal: value.I32(0), expS: &ValueStack{values: []value.Value{}, sp: 0}},
		{s: &ValueStack{values: []value.Value{value.I32(0), value.F32(1.1)}, sp: 0}, expVal: value.F32(1.1), expS: &ValueStack{values: []value.Value{value.I32(0)}, sp: 0}},
	} {
		v, err := d.s.Pop()
		require.NoError(t, err)
		assert.Equal(t, d.expVal, v)
		assert.Equal(t, d.expS, d.s)
	}
}

func TestValueStackPop_Err(t *testing.T) {
	for _, d := range []struct {
		s *ValueStack
		l int
	}{
		{s: &ValueStack{values: []value.Value{value.I32(0)}, sp: 0}, l: 1},
		{s: &ValueStack{values: []value.Value{value.I32(0), value.F32(1.1)}, sp: 0}, l: 2},
	} {
		for i := 0; i < d.l; i++ {
			_, err := d.s.Pop()
			require.NoError(t, err)
		}
		_, err := d.s.Pop()
		require.Error(t, err)
	}
}
