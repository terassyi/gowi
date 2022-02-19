package section

type Section interface {
	Code() SectionCode
}

type SectionCode uint8

const (
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

func (code SectionCode) String() string {
	switch code {
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
