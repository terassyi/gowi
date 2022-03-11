package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type element struct {
	entries []*elementEntry
}

type elementEntry struct {
	index  uint32
	offset []byte // init_expr
	number uint32
	elems  []uint32
}

func newElement(payload []byte) (*element, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewElement: decode count: %w", err)
	}
	entries := make([]*elementEntry, 0, uint32(count))
	for i := 0; i < int(count); i++ {
		index, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewElement: decode index: %w", err)
		}
		offset, err := buf.ReadBytes(END)
		if err != nil {
			return nil, fmt.Errorf("NewElement: decode offset: %w", err)
		}
		number, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewElement: decode number: %w", err)
		}
		elems := make([]uint32, 0, int(number))
		for j := 0; j < int(number); j++ {
			elem, _, err := types.DecodeVarUint32(buf)
			if err != nil {
				return nil, fmt.Errorf("NewElement: decode elems: %w", err)
			}
			elems = append(elems, uint32(elem))
		}
		entries = append(entries, &elementEntry{
			index:  uint32(index),
			offset: offset[:len(offset)-1],
			number: uint32(number),
			elems:  elems,
		})
	}
	return &element{
		entries: entries,
	}, nil
}

func (e *element) detail() (string, error) {
	str := fmt.Sprintf("Element[%d]:\n", len(e.entries))
	for i := 0; i < len(e.entries); i++ {
		switch e.entries[i].offset[0] {
		case 0x41:
			init := "i32"
			iv, _, err := types.DecodeVarInt32(bytes.NewBuffer(e.entries[i].offset[1:]))
			if err != nil {
				return str, err
			}
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%d\n", i, 0, e.entries[i].index, len(e.entries[i].elems), init, iv)
		case 0x42:
			init := "i64"
			iv, _, err := types.DecodeVarInt64(bytes.NewBuffer(e.entries[i].offset[1:]))
			if err != nil {
				return str, err
			}
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%d\n", i, 0, e.entries[i].index, len(e.entries[i].elems), init, iv)
		case 0x43:
			init := "f32"
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%v\n", i, 0, e.entries[i].index, len(e.entries[i].elems), init, e.entries[i].offset[1:])
		case 0x44:
			init := "f64"
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%v\n", i, 0, e.entries[i].index, len(e.entries[i].elems), init, e.entries[i].offset[1:])
		}
		for j := 0; j < len(e.entries[i].elems); j++ {
			str += fmt.Sprintf("  - elem[%d] = func[%d]\n", j, e.entries[i].elems[j])
		}
	}
	return str, nil
}
