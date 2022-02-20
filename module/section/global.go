package section

type Global struct {
	count   uint32
	globals any // TODO global_variable*
}

func NewGlobal(payload []byte) (*Global, error) {
	return &Global{}, nil
}

func (*Global) Code() SectionCode {
	return GLOBAL
}
