package section

type Data struct {
	count   uint32
	entries []DataSegment
}

type DataSegment struct {
	index  uint32
	offset int32
	size   uint32
	data   []byte
}

func (*Data) Code() SectionCode {
	return DATA
}
