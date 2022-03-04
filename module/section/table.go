package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Table struct {
	// count   uint32
	Entries []*types.TableType
}

func NewTable(payload []byte) (*Table, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewTable: decode count: %w", err)
	}
	entries := make([]*types.TableType, 0, int(count))
	for i := 0; i < int(count); i++ {
		entry, err := types.NewTableType(buf)
		if err != nil {
			return nil, fmt.Errorf("NewTable: decode entry: %w", err)
		}
		entries = append(entries, entry)
	}
	return &Table{
		Entries: entries,
	}, nil
}

func (*Table) Code() SectionCode {
	return TABLE
}

func (t *Table) Detail() string {
	str := fmt.Sprintf("%s[%d]:\n", t.Code(), len(t.Entries))
	for i := 0; i < len(t.Entries); i++ {
		str += fmt.Sprintf(" - type[%d] type=%s initial=%d\n", i, t.Entries[i].ElementType, t.Entries[i].Limits.Initial)
	}
	return str
}
