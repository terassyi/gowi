package section

type Memory struct {
	count   uint32
	entries any // TODO memory_type*
}

func (*Memory) Code() SectionCode {
	return MEMORY
}
