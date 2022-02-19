package section

type Start struct {
	index uint32
}

func (*Start) Code() SectionCode {
	return START
}
