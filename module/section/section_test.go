package section

import (
	"errors"
	"testing"
)

func TestSectionNew(t *testing.T) {
	sec, err := New(0x1, []byte{})
	if err != nil {
		t.Fatal(err)
	}
	if sec.Code() != TYPE {
		t.Errorf("want %v, actual: %v", TYPE.String(), sec.Code().String())
	}
}
func TestSectionNew_InvalidSectionCode(t *testing.T) {
	_, err := New(0xf, []byte{})
	if !errors.Is(err, InvalidSectionCode) {
		t.Error(err)
	}
}
