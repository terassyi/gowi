package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/decoder"
	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/debugger"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/stack"
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
	"github.com/terassyi/gowi/validator"
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

func TestInvoke_Numeric(t *testing.T) {
	for _, d := range []struct {
		path   string
		export string
		args   []value.Value
		exp    []value.Value
	}{
		{path: "../examples/const.wasm", export: "i32", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/const.wasm", export: "i64", args: []value.Value{}, exp: []value.Value{value.I64(0x1ff)}},
		// {path: "../examples/const.wasm", export: "f32", args: []value.Value{}, exp: []value.Value{value.F32(0.1)}},
		// {path: "../examples/const.wasm", export: "f64", args: []value.Value{}, exp: []value.Value{value.F64(0.1)}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.NewI32(int32(-2))}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.I32(0x3fffffff), value.I32(1)}, exp: []value.Value{value.I32(0x40000000)}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.I32(0x7fffffff), value.I32(1)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "add", args: []value.Value{value.I32(0x80000000), value.I32(0x80000000)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.NewI32(int32(-2))}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.I32(0x3fffffff), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x40000000)}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.I32(0x7fffffff), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "sub", args: []value.Value{value.I32(0x80000000), value.NewI32(int32(1))}, exp: []value.Value{value.I32(0x7fffffff)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(0x10000000), value.I32(4096)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(0x80000000), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(0x80000000), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(0x7fffffff), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x80000001)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(0x01234567), value.I32(0x76543210)}, exp: []value.Value{value.I32(0x358e7470)}},
		{path: "../examples/i32.wasm", export: "mul", args: []value.Value{value.I32(0x7fffffff), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(0), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(5), value.I32(2)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.NewI32(int32(-5)), value.I32(2)}, exp: []value.Value{value.NewI32(int32(-2))}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(5), value.NewI32(int32(-2))}, exp: []value.Value{value.NewI32(int32(-2))}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.NewI32(int32(-5)), value.NewI32(int32(-2))}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(17), value.I32(7)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(0x80000000), value.I32(2)}, exp: []value.Value{value.I32(0xc0000000)}},
		{path: "../examples/i32.wasm", export: "div_s", args: []value.Value{value.I32(0x80000001), value.I32(1000)}, exp: []value.Value{value.I32(0xffdf3b65)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(0), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(5), value.I32(2)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.NewI32(int32(-5)), value.I32(2)}, exp: []value.Value{value.I32(0x7ffffffd)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(5), value.NewI32(int32(-2))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.NewI32(int32(-5)), value.NewI32(int32(-2))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(17), value.I32(7)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(0x80000000), value.I32(2)}, exp: []value.Value{value.I32(0x40000000)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(0x80000000), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(0x8ff00ff0), value.I32(0x10001)}, exp: []value.Value{value.I32(0x8fef)}},
		{path: "../examples/i32.wasm", export: "div_u", args: []value.Value{value.I32(0x80000001), value.I32(1000)}, exp: []value.Value{value.I32(0x20c49b)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(0x7fffffff), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(0), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(5), value.I32(2)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.NewI32(int32(-5)), value.I32(2)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(7), value.I32(3)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.NewI32(int32(-7)), value.NewI32(int32(-3))}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(17), value.I32(7)}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(0x80000000), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(0x80000000), value.I32(2)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_s", args: []value.Value{value.I32(0x80000001), value.I32(1000)}, exp: []value.Value{value.NewI32(int32(-647))}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(0), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(5), value.I32(2)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(5), value.NewI32(int32(-2))}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.NewI32(int32(-5)), value.I32(2)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.NewI32(int32(-5)), value.NewI32(int32(-2))}, exp: []value.Value{value.NewI32(int32(-5))}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(7), value.I32(3)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(11), value.I32(5)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(17), value.I32(7)}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(0x80000000), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(0x80000000), value.NewI32(int32(2))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(0x8ff00ff0), value.I32(0x10001)}, exp: []value.Value{value.I32(0x8001)}},
		{path: "../examples/i32.wasm", export: "rem_u", args: []value.Value{value.I32(0x80000001), value.I32(1000)}, exp: []value.Value{value.I32(649)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(0x7fffffff), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x7fffffff)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(0x7fffffff), value.I32(0x80000000)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(0xf0f0ffff), value.I32(0xfffff0f0)}, exp: []value.Value{value.I32(0xf0f0f0f0)}},
		{path: "../examples/i32.wasm", export: "and", args: []value.Value{value.I32(0xffffffff), value.I32(0xffffffff)}, exp: []value.Value{value.I32(0xffffffff)}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(0x7fffffff), value.NewI32(int32(-1))}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(0x7fffffff), value.I32(0x80000000)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(0x80000000), value.I32(0)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(0xf0f0ffff), value.I32(0xfffff0f0)}, exp: []value.Value{value.I32(0xffffffff)}},
		{path: "../examples/i32.wasm", export: "or", args: []value.Value{value.I32(0xffffffff), value.I32(0xffffffff)}, exp: []value.Value{value.I32(0xffffffff)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(0x7fffffff), value.I32(0x80000000)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(0x80000000), value.I32(0)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.NewI32(int32(-1)), value.I32(0x80000000)}, exp: []value.Value{value.I32(0x7fffffff)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.NewI32(int32(-1)), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(0xf0f0ffff), value.I32(0xfffff0f0)}, exp: []value.Value{value.I32(0x0f0f0f0f)}},
		{path: "../examples/i32.wasm", export: "xor", args: []value.Value{value.I32(0xffffffff), value.I32(0xffffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(0x7fffffff), value.I32(1)}, exp: []value.Value{value.I32(0xfffffffe)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(0xffffffff), value.I32(1)}, exp: []value.Value{value.I32(0xfffffffe)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(0x80000000), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(0x40000000), value.I32(1)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.I32(31)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.I32(32)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.I32(33)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "shl", args: []value.Value{value.I32(1), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(0x7fffffff)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(0x7fffffff), value.I32(1)}, exp: []value.Value{value.I32(0x3fffffff)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(0x80000000), value.I32(1)}, exp: []value.Value{value.I32(0x40000000)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(0x40000000), value.I32(1)}, exp: []value.Value{value.I32(0x20000000)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.I32(32)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.I32(33)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.I32(1), value.I32(0x80000000)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(32)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(33)}, exp: []value.Value{value.I32(0x7fffffff)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(0x80000000)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(0x7fffffff), value.I32(1)}, exp: []value.Value{value.I32(0x3fffffff)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(0x80000000), value.I32(1)}, exp: []value.Value{value.I32(0xc0000000)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(0x40000000), value.I32(1)}, exp: []value.Value{value.I32(0x20000000)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.I32(32)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.I32(33)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.I32(1), value.I32(0x80000000)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(32)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(33)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(0x7fffffff)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "shr_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(0x80000000)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(1), value.I32(32)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0xabcd9876), value.I32(1)}, exp: []value.Value{value.I32(0x579b30ed)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0xfe00dc00), value.I32(4)}, exp: []value.Value{value.I32(0xe00dc00f)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0xb0c1d2e3), value.I32(5)}, exp: []value.Value{value.I32(0x183a5c76)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0x00008000), value.I32(37)}, exp: []value.Value{value.I32(0x00100000)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0xb0c1d2e3), value.I32(0xff05)}, exp: []value.Value{value.I32(0x183a5c76)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0x769abcdf), value.I32(0xffffffed)}, exp: []value.Value{value.I32(0x579beed3)}},
		{path: "../examples/i32.wasm", export: "rotl", args: []value.Value{value.I32(0x769abcdf), value.I32(0x8000000d)}, exp: []value.Value{value.I32(0x579beed3)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0x80000000)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(1), value.I32(32)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0xff00cc00), value.I32(1)}, exp: []value.Value{value.I32(0x7f806600)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0x00080000), value.I32(4)}, exp: []value.Value{value.I32(0x00008000)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0xb0c1d2e3), value.I32(5)}, exp: []value.Value{value.I32(0x1d860e97)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0x00008000), value.I32(37)}, exp: []value.Value{value.I32(0x00000400)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0xb0c1d2e3), value.I32(0xff05)}, exp: []value.Value{value.I32(0x1d860e97)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0x769abcdf), value.I32(0xffffffed)}, exp: []value.Value{value.I32(0xe6fbb4d5)}},
		{path: "../examples/i32.wasm", export: "rotr", args: []value.Value{value.I32(0x769abcdf), value.I32(0x8000000d)}, exp: []value.Value{value.I32(0xe6fbb4d5)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.I32(0x80000000), value.I32(0x80000000)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.I32(0x7fffffff), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.I32(0x7fffffff), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eq", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.I32(0x80000000), value.I32(0x80000000)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.I32(0x7fffffff), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.NewI32(int32(-1)), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.I32(0x7fffffff), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ne", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "lt_s", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "lt_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "lt_s", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "lt_u", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "lt_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "lt_u", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "gt_s", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "gt_s", args: []value.Value{value.I32(1), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "gt_s", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "gt_u", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "gt_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "gt_u", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "le_s", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "le_s", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "le_s", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "le_u", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "le_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "le_u", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ge_s", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "ge_s", args: []value.Value{value.I32(1), value.NewI32(int32(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "ge_s", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ge_u", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "ge_u", args: []value.Value{value.NewI32(int32(-1)), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "ge_u", args: []value.Value{value.I32(0x80000000), value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eqz", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "eqz", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "eqz", args: []value.Value{value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(0xffffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(32)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(0x00008000)}, exp: []value.Value{value.I32(16)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(0xff)}, exp: []value.Value{value.I32(24)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(0x80000000)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(31)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(2)}, exp: []value.Value{value.I32(30)}},
		{path: "../examples/i32.wasm", export: "clz", args: []value.Value{value.I32(0x7fffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "ctz", args: []value.Value{value.NewI32(int32(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "ctz", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(32)}},
		{path: "../examples/i32.wasm", export: "ctz", args: []value.Value{value.I32(0x00008000)}, exp: []value.Value{value.I32(15)}},
		{path: "../examples/i32.wasm", export: "ctz", args: []value.Value{value.I32(0x00010000)}, exp: []value.Value{value.I32(16)}},
		{path: "../examples/i32.wasm", export: "ctz", args: []value.Value{value.I32(0x80000000)}, exp: []value.Value{value.I32(31)}},
		{path: "../examples/i32.wasm", export: "ctz", args: []value.Value{value.I32(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.NewI32(int32(-1))}, exp: []value.Value{value.I32(32)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0x00008000)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0x80008000)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0x7fffffff)}, exp: []value.Value{value.I32(31)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0xAAAAAAAA)}, exp: []value.Value{value.I32(16)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0x55555555)}, exp: []value.Value{value.I32(16)}},
		{path: "../examples/i32.wasm", export: "popcnt", args: []value.Value{value.I32(0xDEADBEEF)}, exp: []value.Value{value.I32(24)}},
		// i64
		{path: "../examples/i64.wasm", export: "add", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "add", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "add", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.NewI64(int64(-2))}},
		{path: "../examples/i64.wasm", export: "add", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(1))}, exp: []value.Value{value.NewI64(int64(0))}},
		{path: "../examples/i64.wasm", export: "add", args: []value.Value{value.I64(0x3fffffff), value.I64(1)}, exp: []value.Value{value.I64(0x40000000)}},
		{path: "../examples/i64.wasm", export: "add", args: []value.Value{value.NewI64(int64(0x7fffffffffffffff)), value.NewI64(int64(1))}, exp: []value.Value{value.NewI64(uint64(0x8000000000000000))}},
		{path: "../examples/i64.wasm", export: "sub", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "sub", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "sub", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "sub", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.NewI64(int64(-2))}},
		{path: "../examples/i64.wasm", export: "sub", args: []value.Value{value.I64(0x3fffffff), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0x40000000)}},
		{path: "../examples/i64.wasm", export: "sub", args: []value.Value{value.I64(0x7fffffffffffffff), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(1))}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(0x1000000000000000), value.I64(4096)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(0x8000000000000000), value.I64(0)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(0x8000000000000000), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(0x7fffffffffffffff), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0x8000000000000001)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(0x0123456789abcdef), value.I64(0xfedcba9876543210)}, exp: []value.Value{value.I64(0x2236d88fe5618cf0)}},
		{path: "../examples/i64.wasm", export: "mul", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(0), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(5), value.I64(2)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.NewI64(int64(-5)), value.I64(2)}, exp: []value.Value{value.NewI64(int64(-2))}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(5), value.NewI64(int64(-2))}, exp: []value.Value{value.NewI64(int64(-2))}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.NewI64(int64(-5)), value.NewI64(int64(-2))}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(17), value.I64(7)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(0x8000000000000000), value.I64(2)}, exp: []value.Value{value.I64(0xc000000000000000)}},
		{path: "../examples/i64.wasm", export: "div_s", args: []value.Value{value.I64(0x8000000000000001), value.I64(1000)}, exp: []value.Value{value.I64(0xffdf3b645a1cac09)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(0), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(5), value.I64(2)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.NewI64(int64(-5)), value.I64(2)}, exp: []value.Value{value.I64(0x7ffffffffffffffd)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(5), value.NewI64(int64(-2))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.NewI64(int64(-5)), value.NewI64(int64(-2))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(17), value.I64(7)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(0x8000000000000000), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(0x8000000000000000), value.NewI64(int64(2))}, exp: []value.Value{value.I64(0x4000000000000000)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(0x8ff00ff00ff00ff0), value.I64(0x100000001)}, exp: []value.Value{value.I64(0x8ff00fef)}},
		{path: "../examples/i64.wasm", export: "div_u", args: []value.Value{value.I64(0x8000000000000001), value.I64(1000)}, exp: []value.Value{value.I64(0x20c49ba5e353f7)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(0x7fffffffffffffff), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(0), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(5), value.I64(2)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.NewI64(int64(-5)), value.I64(2)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(7), value.I64(3)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.NewI64(int64(-7)), value.NewI64(int64(-3))}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(17), value.I64(7)}, exp: []value.Value{value.I64(3)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(0x8000000000000000), value.I64(2)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_s", args: []value.Value{value.I64(0x8000000000000001), value.I64(1000)}, exp: []value.Value{value.NewI64(int64(-807))}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(0), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(5), value.I64(2)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(5), value.NewI64(int64(-2))}, exp: []value.Value{value.I64(5)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.NewI64(int64(-5)), value.I64(2)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.NewI64(int64(-5)), value.NewI64(int64(-2))}, exp: []value.Value{value.NewI64(int64(-5))}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(7), value.I64(3)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(11), value.I64(5)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(17), value.I64(7)}, exp: []value.Value{value.I64(3)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(0x8000000000000000), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(0x8000000000000000), value.I64(2)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(0x8ff00ff00ff00ff0), value.I64(0x100000001)}, exp: []value.Value{value.I64(0x80000001)}},
		{path: "../examples/i64.wasm", export: "rem_u", args: []value.Value{value.I64(0x8000000000000001), value.I64(1000)}, exp: []value.Value{value.I64(809)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(0x7fffffffffffffff), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0x7fffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(0x8000000000000000)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(0xf0f0ffff), value.I64(0xfffff0f0)}, exp: []value.Value{value.I64(0xf0f0f0f0)}},
		{path: "../examples/i64.wasm", export: "and", args: []value.Value{value.I64(0xffffffffffffffff), value.I64(0xffffffffffffffff)}, exp: []value.Value{value.I64(0xffffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(0x8000000000000000)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(0xf0f0ffff), value.I64(0xfffff0f0)}, exp: []value.Value{value.I64(0xffffffff)}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(0x8000000000000000), value.I64(0)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "or", args: []value.Value{value.I64(0xffffffffffffffff), value.I64(0xffffffffffffffff)}, exp: []value.Value{value.I64(0xffffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(0), value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(0), value.I64(0)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(0x8000000000000000)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(0x8000000000000000), value.I64(0)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.NewI64(int64(-1)), value.I64(0x8000000000000000)}, exp: []value.Value{value.I64(0x7fffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.NewI64(int64(-1)), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "xor", args: []value.Value{value.I64(0xffffffffffffffff), value.I64(0xffffffffffffffff)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(1)}, exp: []value.Value{value.I64(0xfffffffffffffffe)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(0xffffffffffffffff), value.I64(1)}, exp: []value.Value{value.I64(0xfffffffffffffffe)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(0x8000000000000000), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(0x4000000000000000), value.I64(1)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(63)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(64)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(65)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "shl", args: []value.Value{value.I64(1), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.I64(0x7fffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(1)}, exp: []value.Value{value.I64(0x3fffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(0x8000000000000000), value.I64(1)}, exp: []value.Value{value.I64(0x4000000000000000)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(0x4000000000000000), value.I64(1)}, exp: []value.Value{value.I64(0x2000000000000000)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.I64(64)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.I64(65)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(1), value.I64(0x8000000000000000)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.I64(0x8000000000000000), value.I64(63)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.NewI64(int64(-1)), value.I64(64)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.NewI64(int64(-1)), value.I64(65)}, exp: []value.Value{value.I64(0x7fffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.NewI64(int64(-1)), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_u", args: []value.Value{value.NewI64(int64(-1)), value.I64(0x8000000000000000)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(1)}, exp: []value.Value{value.I64(0x3fffffffffffffff)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(0x8000000000000000), value.I64(1)}, exp: []value.Value{value.I64(0xc000000000000000)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(0x4000000000000000), value.I64(1)}, exp: []value.Value{value.I64(0x2000000000000000)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.I64(64)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.I64(65)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(1), value.I64(0x8000000000000000)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.I64(0x8000000000000000), value.I64(63)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.NewI64(int64(-1)), value.I64(64)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.NewI64(int64(-1)), value.I64(65)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.NewI64(int64(-1)), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "shr_s", args: []value.Value{value.NewI64(int64(-1)), value.I64(0x8000000000000000)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(1), value.I64(64)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xabcd987602468ace), value.I64(1)}, exp: []value.Value{value.I64(0x579b30ec048d159d)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xfe000000dc000000), value.I64(4)}, exp: []value.Value{value.I64(0xe000000dc000000f)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xabcd1234ef567809), value.I64(53)}, exp: []value.Value{value.I64(0x013579a2469deacf)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xabd1234ef567809c), value.I64(63)}, exp: []value.Value{value.I64(0x55e891a77ab3c04e)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xabcd1234ef567809), value.I64(0xf5)}, exp: []value.Value{value.I64(0x013579a2469deacf)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xabcd7294ef567809), value.I64(0xffffffffffffffed)}, exp: []value.Value{value.I64(0xcf013579ae529dea)}},
		{path: "../examples/i64.wasm", export: "rotl", args: []value.Value{value.I64(0xabd1234ef567809c), value.I64(0x800000000000003f)}, exp: []value.Value{value.I64(0x55e891a77ab3c04e)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I64(0x8000000000000000)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(1), value.I64(64)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xabcd987602468ace), value.I64(1)}, exp: []value.Value{value.I64(0x55e6cc3b01234567)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xfe000000dc000000), value.I64(4)}, exp: []value.Value{value.I64(0x0fe000000dc00000)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xabcd1234ef567809), value.I64(53)}, exp: []value.Value{value.I64(0x6891a77ab3c04d5e)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xabd1234ef567809c), value.I64(63)}, exp: []value.Value{value.I64(0x57a2469deacf0139)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xabcd1234ef567809), value.I64(0xf5)}, exp: []value.Value{value.I64(0x6891a77ab3c04d5e)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xabcd7294ef567809), value.I64(0xffffffffffffffed)}, exp: []value.Value{value.I64(0x94a77ab3c04d5e6b)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0xabd1234ef567809c), value.I64(0x800000000000003f)}, exp: []value.Value{value.I64(0x57a2469deacf0139)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(1), value.I64(63)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/i64.wasm", export: "rotr", args: []value.Value{value.I64(0x8000000000000000), value.I64(63)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.I64(0), value.I64(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.I64(0x8000000000000000), value.I64(0x8000000000000000)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eq", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.I64(0), value.I64(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.I64(1), value.I64(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.NewI64(int64(-1)), value.I64(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.I64(0x8000000000000000), value.I64(0x8000000000000000)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.I64(0x7fffffffffffffff), value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.NewI64(int64(-1)), value.NewI64(int64(-1))}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "ne", args: []value.Value{value.I64(1), value.I64(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eqz", args: []value.Value{value.I64(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/i64.wasm", export: "eqz", args: []value.Value{value.I64(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "eqz", args: []value.Value{value.I64(0x7fffffff)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(0xffffffffffffffff)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(0)}, exp: []value.Value{value.I64(64)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(0x00008000)}, exp: []value.Value{value.I64(48)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(0xff)}, exp: []value.Value{value.I64(56)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(0x8000000000000000)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(1)}, exp: []value.Value{value.I64(63)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(2)}, exp: []value.Value{value.I64(62)}},
		{path: "../examples/i64.wasm", export: "clz", args: []value.Value{value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "ctz", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "ctz", args: []value.Value{value.I64(0)}, exp: []value.Value{value.I64(64)}},
		{path: "../examples/i64.wasm", export: "ctz", args: []value.Value{value.I64(0x00008000)}, exp: []value.Value{value.I64(15)}},
		{path: "../examples/i64.wasm", export: "ctz", args: []value.Value{value.I64(0x00010000)}, exp: []value.Value{value.I64(16)}},
		{path: "../examples/i64.wasm", export: "ctz", args: []value.Value{value.I64(0x8000000000000000)}, exp: []value.Value{value.I64(63)}},
		{path: "../examples/i64.wasm", export: "ctz", args: []value.Value{value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.I64(64)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0)}, exp: []value.Value{value.I64(0)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0x00008000)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0x8000800080008000)}, exp: []value.Value{value.I64(4)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0x7fffffffffffffff)}, exp: []value.Value{value.I64(63)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0xAAAAAAAA55555555)}, exp: []value.Value{value.I64(32)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0x99999999AAAAAAAA)}, exp: []value.Value{value.I64(32)}},
		{path: "../examples/i64.wasm", export: "popcnt", args: []value.Value{value.I64(0xDEADBEEFDEADBEEF)}, exp: []value.Value{value.I64(48)}},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := instance.New(mod)
		require.NoError(t, err)
		interpreter := &interpreter{
			instance: ins,
			stack:    stack.New(),
			cur:      &current{},
			debubber: debugger.New(debugger.DebugLevelNoLog),
		}
		res, err := interpreter.Invoke(d.export, d.args)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}

func TestInvoke_ControlFlow(t *testing.T) {
	for _, d := range []struct {
		path   string
		export string
		args   []value.Value
		exp    []value.Value
	}{
		// control flow
		{path: "../examples/func1.wasm", export: "add", args: []value.Value{value.I32(0), value.I32(2)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/func1.wasm", export: "add", args: []value.Value{value.I32(13), value.I32(27)}, exp: []value.Value{value.I32(40)}},
		{path: "../examples/call_func1.wasm", export: "getAnswerPlus1", args: []value.Value{}, exp: []value.Value{value.I32(43)}},
		{path: "../examples/call_func_nested.wasm", export: "rootFunc", args: []value.Value{value.I32(1), value.I32(2)}, exp: []value.Value{value.I32(6)}},
		{path: "../examples/block.wasm", export: "singular", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/block.wasm", export: "multi", args: []value.Value{}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/block.wasm", export: "nest", args: []value.Value{}, exp: []value.Value{value.I32(150)}},
		{path: "../examples/block.wasm", export: "deep", args: []value.Value{}, exp: []value.Value{value.I32(150)}},
		{path: "../examples/loop.wasm", export: "singular", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/loop.wasm", export: "multi", args: []value.Value{}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/loop.wasm", export: "nest", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/loop.wasm", export: "deep", args: []value.Value{}, exp: []value.Value{value.I32(150)}},
		{path: "../examples/loop.wasm", export: "param-break", args: []value.Value{}, exp: []value.Value{value.I32(13)}},
		{path: "../examples/loop.wasm", export: "params-break", args: []value.Value{}, exp: []value.Value{value.I32(12)}},
		{path: "../examples/loop.wasm", export: "params-id-break", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/loop.wasm", export: "effects", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/loop.wasm", export: "while", args: []value.Value{value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/loop.wasm", export: "while", args: []value.Value{value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/loop.wasm", export: "while", args: []value.Value{value.I64(2)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/loop.wasm", export: "while", args: []value.Value{value.I64(3)}, exp: []value.Value{value.I64(6)}},
		{path: "../examples/loop.wasm", export: "while", args: []value.Value{value.I64(5)}, exp: []value.Value{value.I64(120)}},
		{path: "../examples/loop.wasm", export: "while", args: []value.Value{value.I64(20)}, exp: []value.Value{value.I64(2432902008176640000)}},
		{path: "../examples/loop.wasm", export: "for", args: []value.Value{value.I64(0)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/loop.wasm", export: "for", args: []value.Value{value.I64(1)}, exp: []value.Value{value.I64(1)}},
		{path: "../examples/loop.wasm", export: "for", args: []value.Value{value.I64(2)}, exp: []value.Value{value.I64(2)}},
		{path: "../examples/loop.wasm", export: "for", args: []value.Value{value.I64(3)}, exp: []value.Value{value.I64(6)}},
		{path: "../examples/loop.wasm", export: "for", args: []value.Value{value.I64(5)}, exp: []value.Value{value.I64(120)}},
		{path: "../examples/loop.wasm", export: "for", args: []value.Value{value.I64(20)}, exp: []value.Value{value.I64(2432902008176640000)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(0), value.I32(7)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(7), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(1), value.I32(2)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(1), value.I32(3)}, exp: []value.Value{value.I32(4)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(1), value.I32(4)}, exp: []value.Value{value.I32(6)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(1), value.I32(100)}, exp: []value.Value{value.I32(2550)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(1), value.I32(101)}, exp: []value.Value{value.I32(2601)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(2), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(2), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(3), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(10), value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(2), value.I32(2)}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/loop.wasm", export: "nesting", args: []value.Value{value.I32(2), value.I32(3)}, exp: []value.Value{value.I32(4)}},
		{path: "../examples/if.wasm", export: "if_func", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/if.wasm", export: "empty", args: []value.Value{value.I32(1)}, exp: []value.Value{}},
		{path: "../examples/if.wasm", export: "empty", args: []value.Value{value.I32(0)}, exp: []value.Value{}},
		{path: "../examples/if.wasm", export: "singular", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/if.wasm", export: "singular", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/if.wasm", export: "multi", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(8), value.I32(1)}},
		{path: "../examples/if.wasm", export: "multi", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(9), value.NewI32(int32(-1))}},
		{path: "../examples/if.wasm", export: "nest", args: []value.Value{value.I32(1), value.I32(1)}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/if.wasm", export: "nest", args: []value.Value{value.I32(1), value.I32(0)}, exp: []value.Value{value.I32(10)}},
		{path: "../examples/if.wasm", export: "nest", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(11)}},
		{path: "../examples/if.wasm", export: "nest", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(10)}},
		{path: "../examples/br.wasm", export: "as-block-first", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br.wasm", export: "as-block-mid", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br.wasm", export: "as-block-last", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br.wasm", export: "as-block-value", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/br.wasm", export: "as-loop-first", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br.wasm", export: "as-loop-mid", args: []value.Value{}, exp: []value.Value{value.I32(4)}},
		{path: "../examples/br.wasm", export: "as-loop-last", args: []value.Value{}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/br.wasm", export: "as-br-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/br.wasm", export: "as-return-value", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/br.wasm", export: "as-return-values", args: []value.Value{}, exp: []value.Value{value.I32(2), value.I32(7)}},
		{path: "../examples/br.wasm", export: "as-if-cond", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/br.wasm", export: "as-if-then", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br.wasm", export: "as-if-else", args: []value.Value{}, exp: []value.Value{value.I32(4)}},
		{path: "../examples/br.wasm", export: "as-select-first", args: []value.Value{}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/br.wasm", export: "as-select-second", args: []value.Value{}, exp: []value.Value{value.I32(6)}},
		{path: "../examples/br.wasm", export: "as-select-cond", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/br.wasm", export: "as-select-all", args: []value.Value{}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/br.wasm", export: "as-call-first", args: []value.Value{}, exp: []value.Value{value.I32(12)}},
		{path: "../examples/br.wasm", export: "as-call-mid", args: []value.Value{}, exp: []value.Value{value.I32(13)}},
		{path: "../examples/br.wasm", export: "as-call-last", args: []value.Value{}, exp: []value.Value{value.I32(14)}},
		{path: "../examples/br.wasm", export: "as-call-all", args: []value.Value{}, exp: []value.Value{value.I32(15)}},
		{path: "../examples/br.wasm", export: "as-local.set-value", args: []value.Value{}, exp: []value.Value{value.I32(17)}},
		{path: "../examples/br.wasm", export: "as-local.tee-value", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/br.wasm", export: "nested-block-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/br.wasm", export: "nested-br-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/br_if.wasm", export: "as-block-first", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br_if.wasm", export: "as-block-mid", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br_if.wasm", export: "as-block-last", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br_if.wasm", export: "as-block-first-value", args: []value.Value{}, exp: []value.Value{value.I32(11)}},
		{path: "../examples/br_if.wasm", export: "as-loop-first", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br_if.wasm", export: "as-loop-mid", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/br_if.wasm", export: "as-loop-last", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br_if.wasm", export: "as-br-value", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/br_if.wasm", export: "as-br_if-cond", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br_if.wasm", export: "as-br_if-value", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/br_if.wasm", export: "as-br_if-value-cond", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/br_if.wasm", export: "as-return-value", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/br_if.wasm", export: "as-if-cond", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/br_if.wasm", export: "as-if-then", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br_if.wasm", export: "as-if-else", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/br_if.wasm", export: "as-select-first", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br_if.wasm", export: "as-select-second", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br_if.wasm", export: "as-select-cond", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/br_if.wasm", export: "as-call-first", args: []value.Value{}, exp: []value.Value{value.I32(12)}},
		{path: "../examples/br_if.wasm", export: "as-call-mid", args: []value.Value{}, exp: []value.Value{value.I32(13)}},
		{path: "../examples/br_if.wasm", export: "as-call-last", args: []value.Value{}, exp: []value.Value{value.I32(14)}},
		{path: "../examples/br_if.wasm", export: "nested-block-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/br_if.wasm", export: "nested-br-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/br_if.wasm", export: "nested-br_if-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/br_if.wasm", export: "nested-br_if-value-cond", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/return.wasm", export: "nullary", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/return.wasm", export: "unary", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/return.wasm", export: "as-func-first", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/return.wasm", export: "as-func-mid", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/return.wasm", export: "as-func-last", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/return.wasm", export: "as-func-value", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/return.wasm", export: "as-block-first", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/return.wasm", export: "as-block-mid", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/return.wasm", export: "as-block-last", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/return.wasm", export: "as-block-value", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/return.wasm", export: "as-loop-first", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/return.wasm", export: "as-loop-mid", args: []value.Value{}, exp: []value.Value{value.I32(4)}},
		{path: "../examples/return.wasm", export: "as-loop-last", args: []value.Value{}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/return.wasm", export: "as-br-value", args: []value.Value{}, exp: []value.Value{value.I32(9)}},
		{path: "../examples/return.wasm", export: "as-return-value", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/return.wasm", export: "as-if-cond", args: []value.Value{}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/return.wasm", export: "as-if-then", args: []value.Value{}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/return.wasm", export: "as-if-else", args: []value.Value{}, exp: []value.Value{value.I32(4)}},
		{path: "../examples/return.wasm", export: "as-select-first", args: []value.Value{}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/return.wasm", export: "as-select-second", args: []value.Value{}, exp: []value.Value{value.I32(6)}},
		{path: "../examples/return.wasm", export: "as-select-cond", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/return.wasm", export: "as-call-first", args: []value.Value{}, exp: []value.Value{value.I32(12)}},
		{path: "../examples/return.wasm", export: "as-call-mid", args: []value.Value{}, exp: []value.Value{value.I32(13)}},
		{path: "../examples/return.wasm", export: "as-call-last", args: []value.Value{}, exp: []value.Value{value.I32(14)}},
		// memory related
		{path: "../examples/load.wasm", export: "as-br-value", args: []value.Value{}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-br_if-cond", args: []value.Value{}, exp: []value.Value{}},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := instance.New(mod)
		require.NoError(t, err)
		interpreter := &interpreter{
			instance: ins,
			stack:    stack.New(),
			cur:      &current{},
			debubber: debugger.New(debugger.DebugLevelNoLog),
		}
		res, err := interpreter.Invoke(d.export, d.args)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}

func TestInvoke_MemoryRelated(t *testing.T) {
	for _, d := range []struct {
		path   string
		export string
		args   []value.Value
		exp    []value.Value
	}{
		// memory related
		{path: "../examples/load.wasm", export: "as-br-value", args: []value.Value{}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-br_if-cond", args: []value.Value{}, exp: []value.Value{}},
		{path: "../examples/load.wasm", export: "as-br_if-value", args: []value.Value{}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-br_if-value-cond", args: []value.Value{}, exp: []value.Value{value.I32(7)}},
		{path: "../examples/load.wasm", export: "as-return-value", args: []value.Value{}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-if-cond", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/load.wasm", export: "as-if-then", args: []value.Value{}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-if-else", args: []value.Value{}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-select-first", args: []value.Value{value.I32(0), value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-select-second", args: []value.Value{value.I32(0), value.I32(0)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/load.wasm", export: "as-select-cond", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/load.wasm", export: "as-call-first", args: []value.Value{}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/load.wasm", export: "as-call-mid", args: []value.Value{}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/load.wasm", export: "as-call-last", args: []value.Value{}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/load.wasm", export: "as-unary-operand", args: []value.Value{}, exp: []value.Value{value.I32(32)}},
		{path: "../examples/load.wasm", export: "as-binary-right", args: []value.Value{}, exp: []value.Value{value.I32(10)}},
		{path: "../examples/load.wasm", export: "as-binary-left", args: []value.Value{}, exp: []value.Value{value.I32(10)}},
		{path: "../examples/load.wasm", export: "as-test-operand", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/load.wasm", export: "as-compare-left", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/load.wasm", export: "as-compare-right", args: []value.Value{}, exp: []value.Value{value.I32(1)}},

		{path: "../examples/memory.wasm", export: "data", args: []value.Value{}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/memory.wasm", export: "i32_load8_s", args: []value.Value{value.NewI32(int32(-1))}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/memory.wasm", export: "i32_load8_s", args: []value.Value{value.I32(100)}, exp: []value.Value{value.I32(100)}},
		{path: "../examples/memory.wasm", export: "i32_load8_s", args: []value.Value{value.I32(0xfedc6543)}, exp: []value.Value{value.I32(0x43)}},
		{path: "../examples/memory.wasm", export: "i32_load8_s", args: []value.Value{value.I32(0x3456cdef)}, exp: []value.Value{value.I32(0xffffffef)}},
		{path: "../examples/memory.wasm", export: "i32_load8_u", args: []value.Value{value.NewI32(int32(-1))}, exp: []value.Value{value.I32(255)}},
		{path: "../examples/memory.wasm", export: "i32_load8_u", args: []value.Value{value.I32(200)}, exp: []value.Value{value.I32(200)}},
		{path: "../examples/memory.wasm", export: "i32_load8_u", args: []value.Value{value.I32(0xfedc6543)}, exp: []value.Value{value.I32(0x43)}},
		{path: "../examples/memory.wasm", export: "i32_load8_u", args: []value.Value{value.I32(0x3456cdef)}, exp: []value.Value{value.I32(0xef)}},
		{path: "../examples/memory.wasm", export: "i32_load16_s", args: []value.Value{value.NewI32(int32(-1))}, exp: []value.Value{value.NewI32(int32(-1))}},
		{path: "../examples/memory.wasm", export: "i32_load16_s", args: []value.Value{value.I32(20000)}, exp: []value.Value{value.I32(20000)}},
		{path: "../examples/memory.wasm", export: "i32_load16_s", args: []value.Value{value.I32(0xfedc6543)}, exp: []value.Value{value.I32(0x6543)}},
		{path: "../examples/memory.wasm", export: "i32_load16_s", args: []value.Value{value.I32(0x3456cdef)}, exp: []value.Value{value.I32(0xffffcdef)}},
		{path: "../examples/memory.wasm", export: "i32_load16_u", args: []value.Value{value.NewI32(int32(-1))}, exp: []value.Value{value.I32(65535)}},
		{path: "../examples/memory.wasm", export: "i32_load16_u", args: []value.Value{value.I32(40000)}, exp: []value.Value{value.I32(40000)}},
		{path: "../examples/memory.wasm", export: "i32_load16_u", args: []value.Value{value.I32(0xfedc6543)}, exp: []value.Value{value.I32(0x6543)}},
		{path: "../examples/memory.wasm", export: "i32_load16_u", args: []value.Value{value.I32(0x3456cdef)}, exp: []value.Value{value.I32(0xcdef)}},
		{path: "../examples/memory.wasm", export: "i64_load8_s", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/memory.wasm", export: "i64_load8_s", args: []value.Value{value.I64(100)}, exp: []value.Value{value.I64(100)}},
		{path: "../examples/memory.wasm", export: "i64_load8_s", args: []value.Value{value.I64(0xfedcba9856346543)}, exp: []value.Value{value.I64(0x43)}},
		{path: "../examples/memory.wasm", export: "i64_load8_s", args: []value.Value{value.I64(0x3456436598bacdef)}, exp: []value.Value{value.I64(0xffffffffffffffef)}},
		{path: "../examples/memory.wasm", export: "i64_load8_u", args: []value.Value{value.I64(200)}, exp: []value.Value{value.I64(200)}},
		{path: "../examples/memory.wasm", export: "i64_load8_u", args: []value.Value{value.I64(0xfedcba9856346543)}, exp: []value.Value{value.I64(0x43)}},
		{path: "../examples/memory.wasm", export: "i64_load8_u", args: []value.Value{value.I64(0x3456436598bacdef)}, exp: []value.Value{value.I64(0xef)}},
		{path: "../examples/memory.wasm", export: "i64_load16_s", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/memory.wasm", export: "i64_load16_s", args: []value.Value{value.I64(20000)}, exp: []value.Value{value.I64(20000)}},
		{path: "../examples/memory.wasm", export: "i64_load16_s", args: []value.Value{value.I64(0xfedcba9856346543)}, exp: []value.Value{value.I64(0x6543)}},
		{path: "../examples/memory.wasm", export: "i64_load16_s", args: []value.Value{value.I64(0x3456436598bacdef)}, exp: []value.Value{value.I64(0xffffffffffffcdef)}},
		{path: "../examples/memory.wasm", export: "i64_load16_u", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.I64(65535)}},
		{path: "../examples/memory.wasm", export: "i64_load16_u", args: []value.Value{value.I64(40000)}, exp: []value.Value{value.I64(40000)}},
		{path: "../examples/memory.wasm", export: "i64_load16_u", args: []value.Value{value.I64(0xfedcba9856346543)}, exp: []value.Value{value.I64(0x6543)}},
		{path: "../examples/memory.wasm", export: "i64_load16_u", args: []value.Value{value.I64(0x3456436598bacdef)}, exp: []value.Value{value.I64(0xcdef)}},
		{path: "../examples/memory.wasm", export: "i64_load32_s", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.NewI64(int64(-1))}},
		{path: "../examples/memory.wasm", export: "i64_load32_s", args: []value.Value{value.I64(20000)}, exp: []value.Value{value.I64(20000)}},
		{path: "../examples/memory.wasm", export: "i64_load32_s", args: []value.Value{value.I64(0xfedcba9856346543)}, exp: []value.Value{value.I64(0x56346543)}},
		{path: "../examples/memory.wasm", export: "i64_load32_s", args: []value.Value{value.I64(0x3456436598bacdef)}, exp: []value.Value{value.I64(0xffffffff98bacdef)}},
		{path: "../examples/memory.wasm", export: "i64_load32_u", args: []value.Value{value.NewI64(int64(-1))}, exp: []value.Value{value.NewI64(int64(4294967295))}},
		{path: "../examples/memory.wasm", export: "i64_load32_u", args: []value.Value{value.I64(40000)}, exp: []value.Value{value.I64(40000)}},
		{path: "../examples/memory.wasm", export: "i64_load32_u", args: []value.Value{value.I64(0xfedcba9856346543)}, exp: []value.Value{value.I64(0x56346543)}},
		{path: "../examples/memory.wasm", export: "i64_load32_u", args: []value.Value{value.I64(0x3456436598bacdef)}, exp: []value.Value{value.I64(0x98bacdef)}},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := instance.New(mod)
		require.NoError(t, err)
		interpreter := &interpreter{
			instance: ins,
			stack:    stack.New(),
			cur:      &current{},
			debubber: debugger.New(debugger.DebugLevelNoLog),
		}
		res, err := interpreter.Invoke(d.export, d.args)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}

func TestInvoke_Recursive(t *testing.T) {
	for _, d := range []struct {
		path   string
		export string
		args   []value.Value
		exp    []value.Value
	}{
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(0)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(2)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(3)}, exp: []value.Value{value.I32(6)}},
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(4)}, exp: []value.Value{value.I32(24)}},
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(5)}, exp: []value.Value{value.I32(120)}},
		{path: "../examples/factorial.wasm", export: "factorial", args: []value.Value{value.I32(10)}, exp: []value.Value{value.I32(3628800)}},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := instance.New(mod)
		require.NoError(t, err)
		interpreter := &interpreter{
			instance: ins,
			stack:    stack.New(),
			cur:      &current{},
			debubber: debugger.New(debugger.DebugLevelNoLog),
		}
		res, err := interpreter.Invoke(d.export, d.args)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}
func TestInvoke_Fibonacci(t *testing.T) {
	for _, d := range []struct {
		path   string
		export string
		args   []value.Value
		exp    []value.Value
	}{
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(2)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(3)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(4)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(5)}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(6)}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(7)}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(8)}, exp: []value.Value{value.I32(13)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(9)}, exp: []value.Value{value.I32(21)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(10)}, exp: []value.Value{value.I32(34)}},
		{path: "../examples/fibonacci.wasm", export: "fib_recursive", args: []value.Value{value.I32(20)}, exp: []value.Value{value.I32(4181)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(1)}, exp: []value.Value{value.I32(0)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(2)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(3)}, exp: []value.Value{value.I32(1)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(4)}, exp: []value.Value{value.I32(2)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(5)}, exp: []value.Value{value.I32(3)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(6)}, exp: []value.Value{value.I32(5)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(7)}, exp: []value.Value{value.I32(8)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(8)}, exp: []value.Value{value.I32(13)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(9)}, exp: []value.Value{value.I32(21)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(10)}, exp: []value.Value{value.I32(34)}},
		{path: "../examples/fibonacci.wasm", export: "fib_iterative", args: []value.Value{value.I32(20)}, exp: []value.Value{value.I32(4181)}},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := instance.New(mod)
		require.NoError(t, err)
		interpreter := &interpreter{
			instance: ins,
			stack:    stack.WithSize(1024, 1024),
			cur:      &current{},
			debubber: debugger.New(debugger.DebugLevelNoLog),
		}
		res, err := interpreter.Invoke(d.export, d.args)
		require.NoError(t, err)
		assert.Equal(t, d.exp, res)
	}
}
func TestStep(t *testing.T) {
	for _, d := range []struct {
		interpreter *interpreter
		instr       instruction.Instruction
		exp         any
		expStack    *stack.Stack
		expCur      *current
	}{
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{{}}), cur: nil},
			instr:       &instruction.I32Const{Imm: 0},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I32(0)}, []stack.Frame{}, []stack.Label{{}}),
			expCur:      nil,
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{{}}), cur: nil},
			instr:       &instruction.I64Const{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I64(1)}, []stack.Frame{}, []stack.Label{{}}),
			expCur:      nil,
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{{}}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.GetLocal{Imm: 0},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I32(0)}, []stack.Frame{}, []stack.Label{{}}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{{}}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.GetLocal{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.F64(0.1)}, []stack.Frame{}, []stack.Label{{}}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.F64(1.5)}, []stack.Frame{}, []stack.Label{}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.SetLocal{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(1.5)}}, label: nil},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.F64(9.0)}, []stack.Frame{}, []stack.Label{}), cur: &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(0.1)}}}},
			instr:       &instruction.TeeLocal{Imm: 1},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.F64(9.0)}, []stack.Frame{}, []stack.Label{}),
			expCur:      &current{frame: &stack.Frame{Module: nil, Locals: []value.Value{value.I32(0), value.F64(9.0)}}, label: nil},
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.I32(0xf0)}, []stack.Frame{}, []stack.Label{}), cur: nil},
			instr:       &instruction.Drop{},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{}, []stack.Frame{}, []stack.Label{}),
			expCur:      nil,
		},
		{
			interpreter: &interpreter{stack: stackWithValueIgnoreError([]value.Value{value.I32(0xff), value.I32(0xee), value.I32(0x0)}, []stack.Frame{}, []stack.Label{}), cur: nil},
			instr:       &instruction.Select{},
			exp:         nil,
			expStack:    stackWithValueIgnoreError([]value.Value{value.I32(0xee)}, []stack.Frame{}, []stack.Label{}),
			expCur:      nil,
		},
	} {
		_, err := d.interpreter.step(d.instr)
		require.NoError(t, err)
		assert.Equal(t, d.expStack, d.interpreter.stack)
		assert.Equal(t, d.expCur, d.interpreter.cur)
	}
}

func stackWithValueIgnoreError(values []value.Value, frames []stack.Frame, labels []stack.Label) *stack.Stack {
	s, _ := stack.WithValue(values, frames, labels)
	return s
}
