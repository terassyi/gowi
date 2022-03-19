package instance

import (
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type Memory struct {
	Type *types.MemoryType
	Data []byte
}

func newMemories(mod *structure.Module) []*Memory {
	mems := make([]*Memory, 0, len(mod.Memories))
	for _, m := range mod.Memories {
		mems = append(mems, &Memory{
			Type: m.Type,
			Data: make([]byte, m.Type.Limits.Min),
		})
	}
	return mems
}

func (*Memory) ExternalValueType() ExternValueType {
	return ExternValTypeMem
}

// https://webassembly.github.io/spec/core/exec/modules.html#growing-memories
func (m *Memory) grow(offset int32, data []byte) error {
	for i, d := range data {
		m.Data[int(offset)+i] = d
	}
	return nil
}
