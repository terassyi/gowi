package section

import "github.com/terassyi/gowi/types"

type Memory struct {
	count   uint32
	entries []types.MemoryType
}

func NewMemory(payload []byte) (*Memory, error) {
	return &Memory{}, nil
}

func (*Memory) Code() SectionCode {
	return MEMORY
}
