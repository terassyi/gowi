package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Export struct {
	// count   uint32
	Entries []*ExportEntry
}

type ExportEntry struct {
	FieldString []byte
	Kind        types.ExternalKind
	Index       uint32
}

func NewExport(payload []byte) (*Export, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewExport: decode count: %w", err)
	}
	entries := make([]*ExportEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		entry := &ExportEntry{}
		fieldLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode fieldLength: %w", err)
		}
		field := make([]byte, int(fieldLength))
		if _, err := buf.Read(field); err != nil {
			return nil, fmt.Errorf("NewExport: decode field_string: %w", err)
		}
		entry.FieldString = field
		externalKind, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode external_kind: %w", err)
		}
		entry.Kind, err = types.NewExternalKind(externalKind)
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode external_kind: %w", err)
		}
		index, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewExport: decode index: %w", err)
		}
		entry.Index = uint32(index)
		entries = append(entries, entry)
	}
	return &Export{
		Entries: entries,
	}, nil
}

func (*Export) Code() SectionCode {
	return EXPORT
}

func (e *Export) Detail() string {
	str := fmt.Sprintf("%s[%d]:\n", e.Code(), len(e.Entries))
	for i := 0; i < len(e.Entries); i++ {
		str += fmt.Sprintf(" - %s[%d] <%s> -> index=%d\n", e.Entries[i].Kind, i, e.Entries[i].FieldString, e.Entries[i].Index)
	}
	return str
}
