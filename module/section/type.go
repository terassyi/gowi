package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Type struct {
	// count   uint32 // count of type entries to follow
	Entries []*types.FuncType
}

func NewType(payload []byte) (*Type, error) {
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
	return &Type{
		Entries: entries,
	}, nil
}

func (*Type) Code() SectionCode {
	return TYPE
}

func (t *Type) Detail() string {
	str := fmt.Sprintf("%s[%d]:\n", t.Code(), len(t.Entries))
	for i := 0; i < len(t.Entries); i++ {
		params := "("
		for j := 0; j < len(t.Entries[i].Params); j++ {
			params += t.Entries[i].Params[j].String()
			if j < len(t.Entries[i].Params)-1 {
				params += ","
			}
		}
		params += ")"
		returns := "("
		for j := 0; j < len(t.Entries[i].Returns); j++ {
			returns += t.Entries[i].Returns[j].String()
			if j < len(t.Entries[i].Returns)-1 {
				returns += ","
			}
		}
		returns += ")"
		str += fmt.Sprintf(" - type[%d] %s -> %s\n", i, params, returns)
	}
	return str
}
