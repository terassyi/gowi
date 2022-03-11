package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type typ struct {
	// count   uint32 // count of type entries to follow
	entries []*types.FuncType
}

func newType(payload []byte) (*typ, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewType: %w", err)
	}
	entries := make([]*types.FuncType, 0, count)
	for i := 0; i < int(count); i++ {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("NewType: %w", err)
		}
		isFunc, err := types.NewValueType(b)
		if err != nil {
			return nil, fmt.Errorf("NewType: %w", err)
		}
		if isFunc != types.FUNC {
			return nil, fmt.Errorf("NewType: func value_type is required: %w", err)
		}
		f, read, err := types.DecodeFuncType(buf.Bytes())
		if err != nil {
			return nil, fmt.Errorf("NewType: %w", err)
		}
		buf.Next(read)
		entries = append(entries, f)
	}
	return &typ{
		entries: entries,
	}, nil
}

func (t *typ) detail() string {
	str := fmt.Sprintf("Type[%d]:\n", len(t.entries))
	for i := 0; i < len(t.entries); i++ {
		params := "("
		for j := 0; j < len(t.entries[i].Params); j++ {
			params += t.entries[i].Params[j].String()
			if j < len(t.entries[i].Params)-1 {
				params += ","
			}
		}
		params += ")"
		returns := "("
		for j := 0; j < len(t.entries[i].Returns); j++ {
			returns += t.entries[i].Returns[j].String()
			if j < len(t.entries[i].Returns)-1 {
				returns += ","
			}
		}
		returns += ")"
		str += fmt.Sprintf(" - type[%d] %s -> %s\n", i, params, returns)
	}
	return str
}
