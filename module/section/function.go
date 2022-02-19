package section

type Function struct {
	count uint32
	types []uint32
}

func (*Function) Code() SectionCode {
	return FUNCTION
}
