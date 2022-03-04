package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Element struct {
	Entries []*ElementEntry
}

type ElementEntry struct {
	Index  uint32
	Offset []byte // init_expr
	Number uint32
	Elems  []uint32
}

func NewElement(payload []byte) (*Element, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewElement: decode count: %w", err)
	}
	entries := make([]*ElementEntry, 0, uint32(count))
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
		entries = append(entries, &ElementEntry{
			Index:  uint32(index),
			Offset: offset[:len(offset)-1],
			Number: uint32(number),
			Elems:  elems,
		})
	}
	return &Element{
		Entries: entries,
	}, nil
}

func (*Element) Code() SectionCode {
	return ELEMENT
}

func (e *Element) Detail() (string, error) {
	str := fmt.Sprintf("%s[%d]:\n", e.Code(), len(e.Entries))
	for i := 0; i < len(e.Entries); i++ {
		switch e.Entries[i].Offset[0] {
		case 0x41:
			init := "i32"
			iv, _, err := types.DecodeVarInt32(bytes.NewBuffer(e.Entries[i].Offset[1:]))
			if err != nil {
				return str, err
			}
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%d\n", i, 0, e.Entries[i].Index, len(e.Entries[i].Elems), init, iv)
		case 0x42:
			init := "i64"
			iv, _, err := types.DecodeVarInt64(bytes.NewBuffer(e.Entries[i].Offset[1:]))
			if err != nil {
				return str, err
			}
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%d\n", i, 0, e.Entries[i].Index, len(e.Entries[i].Elems), init, iv)
		case 0x43:
			init := "f32"
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%v\n", i, 0, e.Entries[i].Index, len(e.Entries[i].Elems), init, e.Entries[i].Offset[1:])
		case 0x44:
			init := "f64"
			str += fmt.Sprintf(" - segment[%d] flags=%d table=%d count=%d - init %s=%v\n", i, 0, e.Entries[i].Index, len(e.Entries[i].Elems), init, e.Entries[i].Offset[1:])
		}
		for j := 0; j < len(e.Entries[i].Elems); j++ {
			str += fmt.Sprintf("  - elem[%d] = func[%d]\n", j, e.Entries[i].Elems[j])
		}
	}
	return str, nil
}
