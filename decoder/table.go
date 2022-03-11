package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type table struct {
	// count   uint32
	entries []*types.TableType
}

func newTable(payload []byte) (*table, error) {
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
	return &table{
		entries: entries,
	}, nil
}

func (t *table) detail() string {
	str := fmt.Sprintf("Table[%d]:\n", len(t.entries))
	for i := 0; i < len(t.entries); i++ {
		str += fmt.Sprintf(" - type[%d] type=%s initial=%d\n", i, t.entries[i].ElementType, t.entries[i].Limits.Min)
	}
	return str
}
