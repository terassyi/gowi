package section

type Type struct {
	count   uint32 // count of type entries to follow
	entries any    // TODO func_type*
}

func (*Type) Code() SectionCode {
	return TYPE
}
