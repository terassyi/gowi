package section

type Element struct {
	count   uint32
	entries []ElementEntry
}

type ElementEntry struct {
	index  uint32
	offset int32
	number uint32
	elems  []uint32
}

func (*Element) Code() SectionCode {
	return ELEMENT
}
