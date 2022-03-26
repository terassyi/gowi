package runtime

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

func TestValidateLocals(t *testing.T) {
	for _, d := range []struct {
		f      *instance.Function
		locals []value.Value
	}{
		{f: &instance.Function{Type: &types.FuncType{Params: []types.ValueType{types.I32, types.F32}, Returns: []types.ValueType{types.F64}}}, locals: []value.Value{value.I32(0), value.F32(1.0)}},
		{f: &instance.Function{Type: &types.FuncType{Params: []types.ValueType{types.I32}, Returns: []types.ValueType{types.F64}}}, locals: []value.Value{value.I32(0)}},
		{f: &instance.Function{Type: &types.FuncType{Params: []types.ValueType{}, Returns: []types.ValueType{}}}, locals: []value.Value{}},
	} {
		err := validateLocals(d.f, d.locals)
		require.NoError(t, err)
	}
}
