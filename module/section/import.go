package section

type Import struct {
	count   uint32
	entries []ImportEntry
}

type ImportEntry struct {
	moduleNameLength uint32
	moduleName       []byte // or stging?
	fieldLength      uint32
	fieldString      []byte
	kind             Kind
}

type Kind uint8 // import kind

func (*Import) Code() SectionCode {
	return IMPORT
}
