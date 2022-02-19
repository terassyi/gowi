package section

type Code struct {
	count  uint32
	bodies any // function_body*
}

func (*Code) Code() SectionCode {
	return CODE
}
