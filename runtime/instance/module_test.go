package instance

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/decoder"
	"github.com/terassyi/gowi/validator"
)

func TestModuleNew(t *testing.T) {
	for _, d := range []struct {
		path string
	}{
		{path: "../../examples/empty_module.wasm"},
		{path: "../../examples/func1.wasm"},
		{path: "../../examples/call_func1.wasm"},
		{path: "../../examples/data1.wasm"},
		{path: "../../examples/elem.wasm"},
		{path: "../../examples/global1.wasm"},
		{path: "../../examples/i32_add.wasm"},
		{path: "../../examples/local1.wasm"},
		{path: "../../examples/mem0.wasm"},
		{path: "../../examples/table.wasm"},
		{path: "../../examples/start0.wasm"},
		{path: "../../examples/import_js.wasm"},
		// I should prepare invalid wasm file to pass test cases
		// {path: "../examples/invalid_table.wasm", res: false},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		_, err = New(mod)
		require.NoError(t, err)
	}
}

func TestModuleGetExports(t *testing.T) {
	for _, d := range []struct {
		path      string
		exportLen int
	}{
		{path: "../../examples/empty_module.wasm", exportLen: 0},
		{path: "../../examples/func1.wasm", exportLen: 1},
		{path: "../../examples/call_func1.wasm", exportLen: 1},
		{path: "../../examples/local1.wasm", exportLen: 8},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := New(mod)
		require.NoError(t, err)
		externals := ins.GetExports()
		assert.Equal(t, d.exportLen, len(externals))
	}
}

func TestModuleGetExport(t *testing.T) {
	for _, d := range []struct {
		path string
		name string
		exp  ExternalValueType
	}{
		{path: "../../examples/func1.wasm", name: "add", exp: ExternalValueTypeFunc},
		{path: "../../examples/call_func1.wasm", name: "getAnswerPlus1", exp: ExternalValueTypeFunc},
		{path: "../../examples/local1.wasm", name: "type-local-i32", exp: ExternalValueTypeFunc},
		{path: "../../examples/start0.wasm", name: "get", exp: ExternalValueTypeFunc},
	} {
		dec, err := decoder.New(d.path)
		require.NoError(t, err)
		mod, err := dec.Decode()
		v, err := validator.New(mod)
		require.NoError(t, err)
		_, err = v.Validate()
		require.NoError(t, err)
		ins, err := New(mod)
		require.NoError(t, err)
		external, err := ins.GetExport(d.name)
		assert.Equal(t, d.exp, external.ExternalValueType())
	}
}
