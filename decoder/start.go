package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type start struct {
	index uint32
}

func newStart(payload []byte) (*start, error) {
	buf := bytes.NewBuffer(payload)
	index, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewStart: decode index: %w", err)
	}
	return &start{
		index: uint32(index),
	}, nil
}

func (s *start) detail() string {
	return fmt.Sprintf("Start:\n- start function: %d", s.index)
}
