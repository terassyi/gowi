package decoder

import (
	"errors"
	"testing"

	"github.com/terassyi/gowi/module"
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
	path := "../examples/func1.wasm"
	data, err := readWasmFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sds, err := decodeSections(data)
	if err != nil {
		t.Fatalf("decodeSections error: %v", err)
	}
	if len(sds) != 4 {
		t.Errorf("want: 4, actual: %d", len(sds))
	}
}
