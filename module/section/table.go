package section

import "github.com/terassyi/gowi/types"

type Table struct {
	count   uint32
	entries []types.TableType
}

func NewTable(payload []byte) (*Table, error) {
	return &Table{}, nil
}

func (*Table) Code() SectionCode {
	return TABLE
}
