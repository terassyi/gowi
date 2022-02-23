package section

import (
	"errors"
	"fmt"
)

type Section interface {
	Code() SectionCode
}

type SectionCode uint8

const (
	CUSTOM   SectionCode = 0x0
	TYPE     SectionCode = 0x1
	IMPORT   SectionCode = 0x2
	FUNCTION SectionCode = 0x3
	TABLE    SectionCode = 0x4
	MEMORY   SectionCode = 0x5
	GLOBAL   SectionCode = 0x6
	EXPORT   SectionCode = 0x7
	START    SectionCode = 0x8
	ELEMENT  SectionCode = 0x9
	CODE     SectionCode = 0xa
	DATA     SectionCode = 0xb
)

var (
	InvalidSectionCode error = errors.New("invalid section code.")
)

func NewSectionCode(val uint8) (SectionCode, error) {
	switch val {
	case uint8(CUSTOM):
		return CUSTOM, nil
	case uint8(TYPE):
		return TYPE, nil
	case uint8(IMPORT):
		return IMPORT, nil
	case uint8(FUNCTION):
		return FUNCTION, nil
	case uint8(TABLE):
		return TABLE, nil
	case uint8(MEMORY):
		return MEMORY, nil
	case uint8(GLOBAL):
		return GLOBAL, nil
	case uint8(EXPORT):
		return EXPORT, nil
	case uint8(START):
		return START, nil
	case uint8(ELEMENT):
		return ELEMENT, nil
	case uint8(CODE):
		return CODE, nil
	case uint8(DATA):
		return DATA, nil
	default:
		return 0xff, fmt.Errorf("%w: %x", InvalidSectionCode, val)
	}
}

func (code SectionCode) String() string {
	switch code {
	case CUSTOM:
		return "CUSTOM"
	case TYPE:
		return "Type"
	case IMPORT:
		return "Import"
	case FUNCTION:
		return "Function"
	case TABLE:
		return "Table"
	case MEMORY:
		return "Memory"
	case GLOBAL:
		return "Global"
	case EXPORT:
		return "Export"
	case START:
		return "Start"
	case ELEMENT:
		return "Element"
	case CODE:
		return "Code"
	case DATA:
		return "Data"
	default:
		return ""
	}
}
