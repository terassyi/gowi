package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Data struct {
	// count   uint32
	Entries []*DataSegment
}

type DataSegment struct {
	Index  uint32
	Offset []byte // init_expr
	Size   uint32
	Data   []byte
}

func NewData(payload []byte) (*Data, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewData: decode count: %w", err)
	}
	entries := make([]*DataSegment, 0, int(count))
	for i := 0; i < int(count); i++ {
		index, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewData: decode index: %w", err)
		}
		offset, err := buf.ReadBytes(END)
		if err != nil {
			return nil, fmt.Errorf("NewData: decode offset: %w", err)
		}
		size, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewData: decode size: %w", err)
		}
		data := make([]byte, int(size))
		if _, err := buf.Read(data); err != nil {
			return nil, fmt.Errorf("NewData: decode data: %w", err)
		}
		entries = append(entries, &DataSegment{
			Index:  uint32(index),
			Offset: offset[:len(offset)-1],
			Size:   uint32(size),
			Data:   data,
		})
	}
	return &Data{
		Entries: entries,
	}, nil
}

func (*Data) Code() SectionCode {
	return DATA
}

func (d *Data) Detail() (string, error) {
	str := fmt.Sprintf("%s[%d]:\n", d.Code(), len(d.Entries))
	for i := 0; i < len(d.Entries); i++ {
		if d.Entries[i].Offset[0] != 0x41 {
			return "", fmt.Errorf("Data Detail: invalid init_expr.")
		}
		init := "i32"
		iv, _, err := types.DecodeVarInt32(bytes.NewBuffer(d.Entries[i].Offset[1:]))
		if err != nil {
			return str, err
		}
		str += fmt.Sprintf(" - segment[%d] memory=%d size=%d init %s=%d\n  - %07x: %v\n", i, d.Entries[i].Index, d.Entries[i].Size, init, iv, iv, d.Entries[i].Data)
	}
	return str, nil
}
