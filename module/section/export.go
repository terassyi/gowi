package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Export struct {
	count   uint32
	entries []*ExportEntry
}

type ExportEntry struct {
	fieldLength uint32
	fieldString []byte
	kind        types.ExternalKind
	index       uint32
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
		entry.fieldLength = uint32(fieldLength)
		entry.fieldString = buf.Bytes()[:fieldLength]
		buf.Next(int(fieldLength))
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
	return &Export{
		count:   uint32(count),
		entries: entries,
	}, nil
}

func (*Export) Code() SectionCode {
	return EXPORT
}
