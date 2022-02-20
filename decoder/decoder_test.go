package decoder

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/terassyi/gowi/module"
	"github.com/terassyi/gowi/module/section"
)

func TestValidateExt_Pass(t *testing.T) {
	pass := "path/to/good.wasm"
	if err := validateExt(pass); err != nil {
		t.Errorf("want: nil, actual: %v", err)
	}
}

func TestValidateExt_Fail(t *testing.T) {
	fail := "path/to/to/ng.wat"
	if err := validateExt(fail); err == nil {
		t.Errorf("want: some error, actual: nil")
	}
}

func TestValidateMajicNumber_Pass(t *testing.T) {
	data := []byte{0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00}
	if err := validateMajicNumber(data); err != nil {
		if errors.Is(err, InvalidMajicNumber) {
			t.Errorf("invalid majic number: want %x", module.MAJIC_NUMBER)
		}
		if errors.Is(err, InvalidWasmVersion) {
			t.Errorf("invalid wasm version: want %x", module.WASM_VERSION)
		}
	}
}

func TestDecodeSections(t *testing.T) {
	for _, f := range []struct {
		path  string
		codes []section.SectionCode
	}{
		{path: "../examples/empty_module.wasm", codes: []section.SectionCode{}},
		{path: "../examples/func1.wasm", codes: []section.SectionCode{section.TYPE, section.FUNCTION, section.EXPORT, section.CODE}},
		{path: "../examples/call_func1.wasm", codes: []section.SectionCode{section.TYPE, section.FUNCTION, section.EXPORT, section.CODE}},
		{path: "../examples/global.wasm", codes: []section.SectionCode{section.TYPE, section.IMPORT, section.FUNCTION, section.EXPORT, section.CODE}},
		{path: "../examples/import_js.wasm", codes: []section.SectionCode{section.TYPE, section.IMPORT, section.FUNCTION, section.EXPORT, section.CODE}},
		{path: "../examples/mem1.wasm", codes: []section.SectionCode{section.TYPE, section.IMPORT, section.FUNCTION, section.EXPORT, section.CODE, section.DATA}},
		{path: "../examples/table.wasm", codes: []section.SectionCode{section.TYPE, section.FUNCTION, section.TABLE, section.EXPORT, section.ELEMENT, section.CODE}},
		{path: "../examples/shared0.wasm", codes: []section.SectionCode{section.TYPE, section.IMPORT, section.FUNCTION, section.ELEMENT, section.CODE}},
		{path: "../examples/shared1.wasm", codes: []section.SectionCode{section.TYPE, section.IMPORT, section.FUNCTION, section.EXPORT, section.CODE}},
	} {
		d, err := readWasmFile(f.path)
		require.NoError(t, err)
		secs, err := decodeSections(d)
		require.NoError(t, err)
		var codes = []section.SectionCode{}
		for _, sec := range secs {
			codes = append(codes, sec.Code())
		}
		assert.Equal(t, f.codes, codes)
		t.Logf("pass %s", f.path)
	}
}
