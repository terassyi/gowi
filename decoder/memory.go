package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type memory struct {
	// count   uint32
	entries []*types.MemoryType
}

func newMemory(payload []byte) (*memory, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewMemory: decode count: %w", err)
	}
	entries := make([]*types.MemoryType, 0, int(count))
	for i := 0; i < int(count); i++ {
		entry, err := types.NewMemoryType(buf)
		if err != nil {
			return nil, fmt.Errorf("NewMemory: decode memory_type: %w", err)
		}
		entries = append(entries, entry)
	}
	return &memory{
		entries: entries,
	}, nil
}

func (m *memory) detail() string {
	str := fmt.Sprintf("Memory[%d]:\n", len(m.entries))
	for i := 0; i < len(m.entries); i++ {
		str += fmt.Sprintf(" - memory[%d] pages: initial=%d\n", i, m.entries[i].Limits.Min)
	}
	return str
}
