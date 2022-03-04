package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Memory struct {
	// count   uint32
	Entries []*types.MemoryType
}

func NewMemory(payload []byte) (*Memory, error) {
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
	return &Memory{
		Entries: entries,
	}, nil
}

func (*Memory) Code() SectionCode {
	return MEMORY
}

func (m *Memory) Detail() string {
	str := fmt.Sprintf("%s[%d]:\n", m.Code(), len(m.Entries))
	for i := 0; i < len(m.Entries); i++ {
		str += fmt.Sprintf(" - memory[%d] pages: initial=%d\n", i, m.Entries[i].Limits.Initial)
	}
	return str
}
