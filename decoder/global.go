package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

const (
	END byte = 0x0b
)

type global struct {
	globals []*globalEntry
}

type globalEntry struct {
	typ  *types.GlobalType
	init []byte
}

func newGlobal(payload []byte) (*global, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewGlobal: decode count: %w", err)
	}
	globals := make([]*globalEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		g, err := newGlobalEntry(buf)
		if err != nil {
			return nil, fmt.Errorf("NewGloabl: decode globals: %w", err)
		}
		globals = append(globals, g)
	}
	return &global{
		globals: globals,
	}, nil
}

func newGlobalEntry(buf *bytes.Buffer) (*globalEntry, error) {
	gt, err := types.NewGloablType(buf)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalEntry: decode global_type: %w", err)
	}
	init, err := buf.ReadBytes(END)
	if err != nil {
		return nil, fmt.Errorf("NewGlobalEntry: decode init_expr: %w", err)
	}
	return &globalEntry{
		typ:  gt,
		init: init[:len(init)-1],
	}, nil
}

func (g *global) detail() (string, error) {
	str := fmt.Sprintf("Global[%d]:\n", len(g.globals))
	for i := 0; i < len(g.globals); i++ {
		mut := 0
		if g.globals[i].typ.Mut {
			mut = 1
		}
		switch g.globals[i].typ.ContentType {
		case types.I32:
			init, _, err := types.DecodeVarInt32(bytes.NewBuffer(g.globals[i].init[1:]))
			if err != nil {
				return "", err
			}
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%d\n", i, g.globals[i].typ.ContentType, mut, g.globals[i].typ.ContentType, init)
		case types.I64:
			init, _, err := types.DecodeVarInt64(bytes.NewBuffer(g.globals[i].init[1:]))
			if err != nil {
				return "", err
			}
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%d\n", i, g.globals[i].typ.ContentType, mut, g.globals[i].typ.ContentType, init)
		case types.F32:
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%v\n", i, g.globals[i].typ.ContentType, mut, g.globals[i].typ.ContentType, g.globals[i].init)
		case types.F64:
			str += fmt.Sprintf(" - global[%d] %s mutable=%d - init %s=%v\n", i, g.globals[i].typ.ContentType, mut, g.globals[i].typ.ContentType, g.globals[i].init)
		}
	}
	return str, nil
}
