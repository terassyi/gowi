package section

type Start struct {
	index uint32
}

func NewStart(payload []byte) (*Start, error) {
	return &Start{}, nil
}

func (*Start) Code() SectionCode {
	return START
}
