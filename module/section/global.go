package section

type Global struct {
	count   uint32
	globals any // TODO global_variable*
}

func (*Global) Code() SectionCode {
	return GLOBAL
}
