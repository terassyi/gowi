package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

func TestTableGrow(t *testing.T) {
	for _, d := range []struct {
		table  *Table
		offset int32
		elems  []uint32
		funcs  []*Function
		exp    []value.Reference
	}{
		{
			table:  &Table{Type: &types.TableType{Limits: &types.Limits{Min: 2}}, Elems: make([]value.Reference, 2)},
			offset: 0,
			elems:  []uint32{0, 1},
			funcs:  []*Function{&Function{}, &Function{}},
			exp:    []value.Reference{&Function{}, &Function{}},
		},
		{
			table:  &Table{Type: &types.TableType{Limits: &types.Limits{Min: 2}}, Elems: make([]value.Reference, 2)},
			offset: 0,
			elems:  []uint32{0, 1, 2},
			funcs:  []*Function{&Function{}, &Function{}, &Function{}},
			exp:    []value.Reference{&Function{}, &Function{}, &Function{}},
		},
		{
			table:  &Table{Type: &types.TableType{Limits: &types.Limits{Min: 2, Max: 3}}, Elems: make([]value.Reference, 2)},
			offset: 0,
			elems:  []uint32{0, 1, 2},
			funcs:  []*Function{&Function{}, &Function{}, &Function{}},
			exp:    []value.Reference{&Function{}, &Function{}, &Function{}},
		},
	} {
		err := d.table.grow(types.ElemTypeFuncref, d.offset, d.elems, d.funcs)
		require.NoError(t, err)
		assert.Equal(t, d.exp, d.table.Elems)
	}
}

func TestTableGrow_Fail(t *testing.T) {
	for _, d := range []struct {
		table  *Table
		offset int32
		elems  []uint32
		funcs  []*Function
		exp    []value.Reference
	}{
		{
			table:  &Table{Type: &types.TableType{Limits: &types.Limits{Min: 2, Max: 3}}, Elems: make([]value.Reference, 2)},
			offset: 0,
			elems:  []uint32{0, 1, 2, 3},
			funcs:  []*Function{&Function{}, &Function{}, &Function{}, &Function{}},
			exp:    []value.Reference{&Function{}, &Function{}, &Function{}},
		},
	} {
		err := d.table.grow(types.ElemTypeFuncref, d.offset, d.elems, d.funcs)
		require.Error(t, err)
	}
}
