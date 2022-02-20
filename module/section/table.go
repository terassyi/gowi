package section

type Table struct {
	count   uint32
	entries any // TODO table_type*
}

func NewTable(payload []byte) (*Table, error) {
	return &Table{}, nil
}

func (*Table) Code() SectionCode {
	return TABLE
}
