package decoder

import (
	"bytes"
	"errors"
	"testing"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	if err := validateMajicNumber(bytes.NewBuffer(data)); err != nil {
		if errors.Is(err, InvalidMajicNumber) {
			t.Errorf("invalid majic number: want %x", MAJIC_NUMBER)
		}
	}
}

func TestDecodeVersion(t *testing.T) {
	data := []byte{0x01, 0x00, 0x00, 0x00}
	ver, err := decodeVersion(bytes.NewBuffer(data))
	require.NoError(t, err)
	assert.Equal(t, uint32(0x01), ver)
}

func TestDecode(t *testing.T) {
	for _, d := range []struct {
		path string
		mod  *mod
	}{
		{
			path: "../examples/empty_module.wasm",
			mod: &mod{
				version: uint32(0x01),
			},
		},
		{
			path: "../examples/func1.wasm",
			mod: &mod{
				version: uint32(0x01),
				typ: &typ{
					entries: []*types.FuncType{
						{
							Params:  []types.ValueType{types.I32, types.I32},
							Returns: []types.ValueType{types.I32},
						},
					},
				},
				function: &function{
					types: []uint32{0x00},
				},
				export: &export{
					entries: []*exportEntry{
						{
							fieldString: []byte{0x61, 0x64, 0x64},
							kind:        types.EXTERNAL_KIND_FUNCTION,
							index:       0x00,
						},
					},
				},
				code: &code{
					bodies: []*functionBody{
						{
							locals: []*localEntry{},
							// Code:   []byte{0x20, 0x00, 0x20, 0x01, 0x6a, 0x0b},
							code: []instruction.Instruction{&instruction.GetLocal{Imm: uint32(0x00)}, &instruction.GetLocal{Imm: uint32(0x01)}, &instruction.I32Add{}, &instruction.End{}},
						},
					},
				},
			},
		},
	} {
		m, err := decode(d.path)
		require.NoError(t, err)
		assert.Equal(t, d.mod.version, m.version)
		if m.typ != nil && d.mod.typ != nil {
			assert.Equal(t, d.mod.typ, m.typ)
		}
		if m.function != nil && d.mod.function != nil {
			assert.Equal(t, d.mod.function, m.function)
		}
		if m.code != nil && d.mod.code != nil {
			assert.Equal(t, d.mod.code, m.code)
		}
	}
}

func TestNewSectionDecoder(t *testing.T) {
	for _, d := range []struct {
		payload []byte
		sd      *sectionDecoder
	}{
		{
			payload: []byte{0x01, 0x07, 0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f},
			sd: &sectionDecoder{
				id:            TYPE,
				payloadLength: uint32(0x07),
				payloadData:   []byte{0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f},
			},
		},
		{
			payload: []byte{0x02, 0x0f, 0x01, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x03, 0x6c, 0x6f, 0x67, 0x00, 0x00},
			sd: &sectionDecoder{
				id:            IMPORT,
				payloadLength: uint32(0x0f),
				payloadData:   []byte{0x01, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x03, 0x6c, 0x6f, 0x67, 0x00, 0x00},
			},
		},
		{
			payload: []byte{0x03, 0x02, 0x01, 0x01},
			sd: &sectionDecoder{
				id:            FUNCTION,
				payloadLength: uint32(0x02),
				payloadData:   []byte{0x01, 0x01},
			},
		},
		{
			payload: []byte{0x07, 0x23, 0x02, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x69, 0x33, 0x32, 0x00, 0x00, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x69, 0x36, 0x34, 0x00, 0x01},
			sd: &sectionDecoder{
				id:            EXPORT,
				payloadLength: uint32(0x23),
				payloadData:   []byte{0x02, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x69, 0x33, 0x32, 0x00, 0x00, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x69, 0x36, 0x34, 0x00, 0x01},
			},
		},
	} {
		sd, err := newSectionDecoder(bytes.NewBuffer(d.payload))
		require.NoError(t, err)
		assert.Equal(t, d.sd, sd)
	}
}
