package section

type Code struct {
	count  uint32
	bodies any // function_body*
}

func NewCode(payload []byte) (*Code, error) {
	return &Code{}, nil
}

func (*Code) Code() SectionCode {
	return CODE
}
