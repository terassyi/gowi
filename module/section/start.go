package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Start struct {
	Index uint32
}

func NewStart(payload []byte) (*Start, error) {
	buf := bytes.NewBuffer(payload)
	index, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewStart: decode index: %w", err)
	}
	return &Start{
		Index: uint32(index),
	}, nil
}

func (*Start) Code() SectionCode {
	return START
}

func (s *Start) Detail() string {
	return fmt.Sprintf("%s:\n- start function: %d", s.Code(), s.Index)
}
