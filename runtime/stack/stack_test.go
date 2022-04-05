package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

func TestStackPushValue(t *testing.T) {
	for _, d := range []struct {
		s   *Stack
		val value.Value
		exp *Stack
	}{
		{s: stackWithValueIgnoreError([]value.Value{}, []Frame{}, []Label{{ValCounter: 0}}), val: value.I32(0), exp: stackWithValueIgnoreError([]value.Value{value.I32(0)}, []Frame{}, []Label{{ValCounter: 1}})},
		{s: stackWithValueIgnoreError([]value.Value{}, []Frame{}, []Label{{ValCounter: 5}}), val: value.I32(0), exp: stackWithValueIgnoreError([]value.Value{value.I32(0)}, []Frame{}, []Label{{ValCounter: 6}})},
	} {
		err := d.s.PushValue(d.val)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.s)
	}
}

func TestStackPopValues(t *testing.T) {
	for _, d := range []struct {
		s   *Stack
		n   int
		exp *Stack
	}{
		{s: stackWithValueIgnoreError([]value.Value{value.I32(0)}, []Frame{}, []Label{{ValCounter: 0}}), n: 0, exp: stackWithValueIgnoreError([]value.Value{value.I32(0)}, []Frame{}, []Label{{ValCounter: 0}})},
		{s: stackWithValueIgnoreError([]value.Value{value.I32(0)}, []Frame{}, []Label{{ValCounter: 5}}), n: 1, exp: stackWithValueIgnoreError([]value.Value{}, []Frame{}, []Label{{ValCounter: 4}})},
		{s: stackWithValueIgnoreError([]value.Value{value.I32(0), value.F32(0.1), value.I32(10), value.F32(6.666)}, []Frame{}, []Label{{ValCounter: 4}}), n: 3, exp: stackWithValueIgnoreError([]value.Value{value.I32(0)}, []Frame{}, []Label{{ValCounter: 1}})},
	} {
		_, err := d.s.PopValues(d.n)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.s)
	}

}
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

func TestValueStackValidate(t *testing.T) {
	for _, d := range []struct {
		stack *ValueStack
		types []types.ValueType
		exp   error
	}{
		{stack: &ValueStack{values: []value.Value{value.I32(0)}}, types: []types.ValueType{types.I32}, exp: nil},
		{stack: &ValueStack{values: []value.Value{value.I32(0), value.I32(1), value.I32(2), value.I32(3)}}, types: []types.ValueType{types.I32, types.I32, types.I32}, exp: nil},
		{stack: &ValueStack{values: []value.Value{value.I32(0), value.I64(1)}}, types: []types.ValueType{types.I64, types.I32}, exp: nil},
		{stack: &ValueStack{values: []value.Value{value.I32(0), value.I64(1)}}, types: []types.ValueType{types.I32, types.I64}, exp: ValueStackTypeNotMatch},
	} {
		err := d.stack.Validate(d.types)
		assert.Equal(t, d.exp, err)
	}
}

func TestValueStackPopNRev(t *testing.T) {
	for _, d := range []struct {
		stack  *ValueStack
		values []value.Value
	}{
		{stack: &ValueStack{values: []value.Value{value.I32(0)}}, values: []value.Value{value.I32(1), value.I32(2)}},
		{stack: &ValueStack{values: []value.Value{value.I32(0), value.I32(1), value.I32(2), value.I32(3)}}, values: []value.Value{value.I32(0xff), value.F32(0.1)}},
		{stack: &ValueStack{values: []value.Value{value.I32(0), value.I64(1)}}, values: []value.Value{}},
	} {
		for _, v := range d.values {
			err := d.stack.Push(v)
			require.NoError(t, err)
		}
		res, err := d.stack.PopNRev(len(d.values))
		require.NoError(t, err)
		assert.Equal(t, d.values, res)
	}
}

func stackWithValueIgnoreError(values []value.Value, frames []Frame, labels []Label) *Stack {
	s, _ := WithValue(values, frames, labels)
	return s
}
