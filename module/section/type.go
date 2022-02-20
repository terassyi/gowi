package section

type Type struct {
	count   uint32 // count of type entries to follow
	entries any    // TODO func_type*
}

func NewType(payload []byte) (*Type, error) {
	return &Type{}, nil
}

func (*Type) Code() SectionCode {
	return TYPE
}
