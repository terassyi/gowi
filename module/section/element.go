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
