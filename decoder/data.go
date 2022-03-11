package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type data struct {
	// count   uint32
	entries []*dataSegment
}

type dataSegment struct {
	index  uint32
	offset []byte // init_expr
	size   uint32
	data   []byte
}

func newData(payload []byte) (*data, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewData: decode count: %w", err)
	}
	entries := make([]*dataSegment, 0, int(count))
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
		entries = append(entries, &dataSegment{
			index:  uint32(index),
			offset: offset[:len(offset)-1],
			size:   uint32(size),
			data:   data,
		})
	}
	return &data{
		entries: entries,
	}, nil
}

func (d *data) detail() (string, error) {
	str := fmt.Sprintf("data[%d]:\n", len(d.entries))
	for i := 0; i < len(d.entries); i++ {
		if d.entries[i].offset[0] != 0x41 {
			return "", fmt.Errorf("Data Detail: invalid init_expr.")
		}
		init := "i32"
		iv, _, err := types.DecodeVarInt32(bytes.NewBuffer(d.entries[i].offset[1:]))
		if err != nil {
			return str, err
		}
		str += fmt.Sprintf(" - segment[%d] memory=%d size=%d init %s=%d\n  - %07x: %v\n", i, d.entries[i].index, d.entries[i].size, init, iv, iv, d.entries[i].data)
	}
	return str, nil
}
