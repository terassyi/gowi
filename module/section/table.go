package section

type Table struct {
	count   uint32
	entries any // TODO table_type*
}

func (*Table) Code() SectionCode {
	return TABLE
}
