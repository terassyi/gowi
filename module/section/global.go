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
	Globals []*GlobalEntry
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
		Globals: globals,
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

func (g *Global) Detail() (string, error) {
	str := fmt.Sprintf("%s[%d]:\n", g.Code(), len(g.Globals))
	for i := 0; i < len(g.Globals); i++ {
		mut := 0
		if g.Globals[i].Type.Mut {
			mut = 1
		}
		switch g.Globals[i].Type.ContentType {
		case types.I32:
			init, _, err := types.DecodeVarInt32(bytes.NewBuffer(g.Globals[i].Init[1:]))
			if err != nil {
				return "", err
			}
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%d\n", i, g.Globals[i].Type.ContentType, mut, g.Globals[i].Type.ContentType, init)
		case types.I64:
			init, _, err := types.DecodeVarInt64(bytes.NewBuffer(g.Globals[i].Init[1:]))
			if err != nil {
				return "", err
			}
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%d\n", i, g.Globals[i].Type.ContentType, mut, g.Globals[i].Type.ContentType, init)
		case types.F32:
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%v\n", i, g.Globals[i].Type.ContentType, mut, g.Globals[i].Type.ContentType, g.Globals[i].Init)
		case types.F64:
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%v\n", i, g.Globals[i].Type.ContentType, mut, g.Globals[i].Type.ContentType, g.Globals[i].Init)
		}
	}
	return str, nil
}
