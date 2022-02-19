package section

type Export struct {
	count   uint32
	entries []ExportEntry
}

type ExportEntry struct {
	fieldLength uint32
	fieldString []byte
	kind        Kind // external_kind
	index       uint32
}

func (*Export) Code() SectionCode {
	return EXPORT
}
