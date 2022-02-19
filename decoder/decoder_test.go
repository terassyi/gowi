package decoder

import "testing"

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
