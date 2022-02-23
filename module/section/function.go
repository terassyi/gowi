package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type Function struct {
	Types []uint32
}

func NewFunction(payload []byte) (*Function, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewFunction: decode count: %w", err)
	}
	typs := make([]uint32, 0, count)
	for i := 0; i < int(count); i++ {
		n, _, err := types.DecodeVarUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("NewFunction: decode types: %w", err)
		}
		typs = append(typs, uint32(n))
	}
	return &Function{
		Types: typs,
	}, nil
}

func (*Function) Code() SectionCode {
	return FUNCTION
}
