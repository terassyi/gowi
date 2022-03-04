package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Import struct {
	Entries []*ImportEntry
}

type ImportEntry struct {
	ModuleNameLength uint32
	ModuleName       []byte // or stging?
	FieldLength      uint32
	FieldString      []byte
	Kind             types.ExternalKind
	Type             any
}

// type ImportEntryType interface {
// 	types.VarUint32 | types.TableType | types.MemoryType
// }

func NewImport(payload []byte) (*Import, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewImport: decode count: %w", err)
	}
	entries := make([]*ImportEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		entry := &ImportEntry{}
		moduleNameLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode module_len: %w", err)
		}
		entry.ModuleNameLength = uint32(moduleNameLength)
		name := make([]byte, int(moduleNameLength))
		if _, err := buf.Read(name); err != nil {
			return nil, fmt.Errorf("NewImport: decode module_name: %w", err)
		}
		entry.ModuleName = name
		fieldLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode field_len: %w", err)
		}
		entry.FieldLength = uint32(fieldLength)
		field := make([]byte, int(fieldLength))
		if _, err := buf.Read(field); err != nil {
			return nil, fmt.Errorf("NewImport: decode field_string: %w", err)
		}
		entry.FieldString = field
		b, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode external_kind: %w", err)
		}
		kind, err := types.NewExternalKind(b)
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode external_kind: %w", err)
		}
		entry.Kind = kind
		switch entry.Kind {
		case types.EXTERNAL_KIND_FUNCTION:
			t, _, err := types.DecodeVarUint32(buf)
			if err != nil {
				return nil, fmt.Errorf("NewImport: decode type: %w", err)
			}
			entry.Type = t
		case types.EXTERNAL_KIND_TABLE:
			t, err := types.NewTableType(buf)
			if err != nil {
				return nil, fmt.Errorf("NewImport: decode type: %w", err)
			}
			// next
			entry.Type = t
		case types.EXTERNAL_KIND_MEMORY:
			t, err := types.NewMemoryType(buf)
			if err != nil {
				return nil, fmt.Errorf("NewImport: decode type: %w", err)
			}
			// next
			entry.Type = t
		}
		entries = append(entries, entry)
	}

	return &Import{
		Entries: entries,
	}, nil
}

func (*Import) Code() SectionCode {
	return IMPORT
}

func (i *Import) Detail() string {
	str := fmt.Sprintf("Import[%d]:\n", len(i.Entries))
	for j := 0; j < len(i.Entries); j++ {
		str += fmt.Sprintf(" - %s <%s.%s>\n", i.Entries[j].Kind, i.Entries[j].ModuleName, i.Entries[j].FieldString)
	}
	return str
}
