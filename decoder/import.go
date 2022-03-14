package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type imports struct {
	entries []*importEntry
}

type importEntry struct {
	moduleNameLength uint32
	moduleName       []byte // or stging?
	fieldLength      uint32
	fieldString      []byte
	kind             types.ExternalKind
	typ              any
}

func fromImportEntry(kind types.ExternalKind, val any) (*structure.ImportDesc, error) {
	switch structure.DescType(kind) {
	case structure.DescTypeFunc:
		return &structure.ImportDesc{
			Type: structure.DescType(kind),
			Func: uint32(val.(types.VarUint32)),
		}, nil
	case structure.DescTypeTable:
		return &structure.ImportDesc{
			Type:  structure.DescType(kind),
			Table: val.(*types.TableType),
		}, nil
	case structure.DescTypeMemory:
		return &structure.ImportDesc{
			Type: structure.DescType(kind),
			Mem:  val.(*types.MemoryType),
		}, nil
	case structure.DescTypeGlobal:
		return &structure.ImportDesc{
			Type:   structure.DescType(kind),
			Global: val.(*types.GlobalType),
		}, nil
	default:
		return nil, fmt.Errorf("fromImportEntry: %w", types.InvalidExternalKind)
	}
}

func newImport(payload []byte) (*imports, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewImport: decode count: %w", err)
	}
	entries := make([]*importEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		entry := &importEntry{}
		moduleNameLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode module_len: %w", err)
		}
		entry.moduleNameLength = uint32(moduleNameLength)
		name := make([]byte, int(moduleNameLength))
		if _, err := buf.Read(name); err != nil {
			return nil, fmt.Errorf("NewImport: decode module_name: %w", err)
		}
		entry.moduleName = name
		fieldLength, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode field_len: %w", err)
		}
		entry.fieldLength = uint32(fieldLength)
		field := make([]byte, int(fieldLength))
		if _, err := buf.Read(field); err != nil {
			return nil, fmt.Errorf("NewImport: decode field_string: %w", err)
		}
		entry.fieldString = field
		b, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode external_kind: %w", err)
		}
		kind, err := types.NewExternalKind(b)
		if err != nil {
			return nil, fmt.Errorf("NewImport: decode external_kind: %w", err)
		}
		entry.kind = kind
		switch entry.kind {
		case types.EXTERNAL_KIND_FUNCTION:
			t, _, err := types.DecodeVarUint32(buf)
			if err != nil {
				return nil, fmt.Errorf("NewImport: decode type: %w", err)
			}
			entry.typ = t
		case types.EXTERNAL_KIND_TABLE:
			t, err := types.NewTableType(buf)
			if err != nil {
				return nil, fmt.Errorf("NewImport: decode type: %w", err)
			}
			// next
			entry.typ = t
		case types.EXTERNAL_KIND_MEMORY:
			t, err := types.NewMemoryType(buf)
			if err != nil {
				return nil, fmt.Errorf("NewImport: decode type: %w", err)
			}
			// next
			entry.typ = t
		}
		entries = append(entries, entry)
	}

	return &imports{
		entries: entries,
	}, nil
}

func (i *imports) detail() string {
	str := fmt.Sprintf("Import[%d]:\n", len(i.entries))
	for j := 0; j < len(i.entries); j++ {
		str += fmt.Sprintf(" - %s <%s.%s>\n", i.entries[j].kind, i.entries[j].moduleName, i.entries[j].fieldString)
	}
	return str
}
