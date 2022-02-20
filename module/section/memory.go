package section

type Memory struct {
	count   uint32
	entries any // TODO memory_type*
}

func NewMemory(payload []byte) (*Memory, error) {
	return &Memory{}, nil
}

func (*Memory) Code() SectionCode {
	return MEMORY
}
