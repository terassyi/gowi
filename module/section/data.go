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

func NewData(payload []byte) (*Data, error) {
	return &Data{}, nil
}

func (*Data) Code() SectionCode {
	return DATA
}
