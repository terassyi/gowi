package instance

import (
	"fmt"

	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

const (
	PAGE_SIZE uint32 = 65536 // 64KB
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
			Data: make([]byte, PAGE_SIZE*m.Type.Limits.Min),
		})
	}
	return mems
}

func (*Memory) ExternalValueType() ExternValueType {
	return ExternValTypeMem
}

// https://webassembly.github.io/spec/core/exec/modules.html#growing-memories
func (m *Memory) initData(offset int32, data []byte) error {
	if int(offset)+len(data) > len(m.Data) {
		return fmt.Errorf("Out of bounds memory access")
	}
	for i, d := range data {
		m.Data[int(offset)+i] = d
	}
	return nil
}
