package section

import "github.com/terassyi/gowi/types"

type Global struct {
	count   uint32
	globals any // TODO global_variable*
}

type GlobalEntry struct {
	Type types.GlobalType
	init uint8
}

func NewGlobal(payload []byte) (*Global, error) {
	return &Global{}, nil
}

func (*Global) Code() SectionCode {
	return GLOBAL
}
