package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type export struct {
	// count   uint32
	entries []*exportEntry
}

type exportEntry struct {
	fieldString []byte
	kind        types.ExternalKind
	index       uint32
}

func newExport(payload []byte) (*export, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewExport: decode count: %w", err)
	}
	entries := make([]*exportEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		entry := &exportEntry{}
		fieldLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode fieldLength: %w", err)
		}
		field := make([]byte, int(fieldLength))
		if _, err := buf.Read(field); err != nil {
			return nil, fmt.Errorf("NewExport: decode field_string: %w", err)
		}
		entry.fieldString = field
		externalKind, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode external_kind: %w", err)
		}
		entry.kind, err = types.NewExternalKind(externalKind)
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode external_kind: %w", err)
		}
		index, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode index: %w", err)
		}
		entry.index = uint32(index)
		entries = append(entries, entry)
	}
	return &export{
		entries: entries,
	}, nil
}

func (e *export) detail() string {
	str := fmt.Sprintf("Export[%d]:\n", len(e.entries))
	for i := 0; i < len(e.entries); i++ {
		str += fmt.Sprintf(" - %s[%d] <%s> -> index=%d\n", e.entries[i].kind, i, e.entries[i].fieldString, e.entries[i].index)
	}
	return str
}
