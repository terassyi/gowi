package section

type Custom struct{}

func NewCustom(payload []byte) (*Custom, error) {
	return &Custom{}, nil
}

func (*Custom) Code() SectionCode {
	return CUSTOM
}

func (*Custom) Detail() string {
	return "not implemented."
}
