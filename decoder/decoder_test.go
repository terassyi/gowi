package decoder

import (
	"bytes"
	"errors"
	"testing"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"

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
	if err := validateMajicNumber(bytes.NewBuffer(data)); err != nil {
		if errors.Is(err, InvalidMajicNumber) {
			t.Errorf("invalid majic number: want %x", module.MAJIC_NUMBER)
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
		mod  *module.Module
	}{
		{
			path: "../examples/empty_module.wasm",
			mod: &module.Module{
				Version: uint32(0x01),
			},
		},
		{
			path: "../examples/func1.wasm",
			mod: &module.Module{
				Version: uint32(0x01),
				Type: &section.Type{
					Entries: []*types.FuncType{
						{
							Params:  []types.ValueType{types.I32, types.I32},
							Returns: []types.ValueType{types.I32},
						},
					},
				},
				Function: &section.Function{
					Types: []uint32{0x00},
				},
				Export: &section.Export{
					Entries: []*section.ExportEntry{
						{
							FieldString: []byte{0x61, 0x64, 0x64},
							Kind:        types.EXTERNAL_KIND_FUNCTION,
							Index:       0x00,
						},
					},
				},
				Code: &section.Code{
					Bodies: []*section.FunctionBody{
						{
							Locals: []*section.LocalEntry{},
							// Code:   []byte{0x20, 0x00, 0x20, 0x01, 0x6a, 0x0b},
							Code: []instruction.Instruction{&instruction.GetLocal{Imm: uint32(0x00)}, &instruction.GetLocal{Imm: uint32(0x01)}, &instruction.I32Add{}, &instruction.End{}},
						},
					},
				},
			},
		},
	} {
		m, err := decode(d.path)
		require.NoError(t, err)
		assert.Equal(t, d.mod.Version, m.Version)
		if m.Type != nil && d.mod.Type != nil {
			assert.Equal(t, d.mod.Type, m.Type)
		}
		if m.Function != nil && d.mod.Function != nil {
			assert.Equal(t, d.mod.Function, m.Function)
		}
		if m.Code != nil && d.mod.Code != nil {
			assert.Equal(t, d.mod.Code, m.Code)
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
				id:            section.TYPE,
				payloadLength: uint32(0x07),
				payloadData:   []byte{0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01, 0x7f},
			},
		},
		{
			payload: []byte{0x02, 0x0f, 0x01, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x03, 0x6c, 0x6f, 0x67, 0x00, 0x00},
			sd: &sectionDecoder{
				id:            section.IMPORT,
				payloadLength: uint32(0x0f),
				payloadData:   []byte{0x01, 0x07, 0x63, 0x6f, 0x6e, 0x73, 0x6f, 0x6c, 0x65, 0x03, 0x6c, 0x6f, 0x67, 0x00, 0x00},
			},
		},
		{
			payload: []byte{0x03, 0x02, 0x01, 0x01},
			sd: &sectionDecoder{
				id:            section.FUNCTION,
				payloadLength: uint32(0x02),
				payloadData:   []byte{0x01, 0x01},
			},
		},
		{
			payload: []byte{0x07, 0x23, 0x02, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x69, 0x33, 0x32, 0x00, 0x00, 0x0e, 0x74, 0x79, 0x70, 0x65, 0x2d, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x2d, 0x69, 0x36, 0x34, 0x00, 0x01},
			sd: &sectionDecoder{
				id:            section.EXPORT,
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
