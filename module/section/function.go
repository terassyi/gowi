package section

type Function struct {
	count uint32
	types []uint32
}

func NewFunction(payload []byte) (*Function, error) {
	return &Function{}, nil
}

func (*Function) Code() SectionCode {
	return FUNCTION
}
