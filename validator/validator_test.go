package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/decoder"
)

func TestValidate(t *testing.T) {
	for _, d := range []struct {
		path string
		res  bool
	}{
		{path: "../examples/empty_module.wasm", res: true},
		{path: "../examples/func1.wasm", res: true},
		{path: "../examples/call_func1.wasm", res: true},
		{path: "../examples/data1.wasm", res: true},
		{path: "../examples/elem.wasm", res: true},
		{path: "../examples/global1.wasm", res: true},
		{path: "../examples/i32_add.wasm", res: true},
		{path: "../examples/local1.wasm", res: true},
		{path: "../examples/mem0.wasm", res: true},
		{path: "../examples/table.wasm", res: true},
		{path: "../examples/start0.wasm", res: true},
		// {path: "../examples/shared0.wasm", res: true},
		// {path: "../examples/shared1.wasm", res: true},
		{path: "../examples/import_js.wasm", res: true},
		// I should prepare invalid wasm file to pass test cases
		// {path: "../examples/invalid_table.wasm", res: false},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		require.NoError(t, err)
		v, err := New(mod)
		require.NoError(t, err)
		res, err := v.Validate()
		require.NoError(t, err)
		assert.Equal(t, d.res, res)
	}
}
