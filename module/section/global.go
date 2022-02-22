package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

const (
	END byte = 0x0b
)

type Global struct {
	count   uint32
	globals []*GlobalEntry
}

type GlobalEntry struct {
	Type *types.GlobalType
	Init []byte
}

func NewGlobal(payload []byte) (*Global, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewGlobal: decode count: %w", err)
	}
	globals := make([]*GlobalEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		g, err := NewGlobalEntry(buf)
		if err != nil {
			return nil, fmt.Errorf("NewGloabl: decode globals: %w", err)
		}
		globals = append(globals, g)
	}
	return &Global{
		count:   uint32(count),
		globals: globals,
	}, nil
}

func NewGlobalEntry(buf *bytes.Buffer) (*GlobalEntry, error) {
	gt, err := types.NewGloablType(buf)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalEntry: decode global_type: %w", err)
	}
	init, err := buf.ReadBytes(END)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalEntry: decode init_expr: %w", err)
	}
	return &GlobalEntry{
		Type: gt,
		Init: init[:len(init)-1],
	}, nil
}

func (*Global) Code() SectionCode {
	return GLOBAL
}
