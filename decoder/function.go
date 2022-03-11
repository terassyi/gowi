package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type function struct {
	types []uint32
}

func newFunction(payload []byte) (*function, error) {
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
	return &function{
		types: typs,
	}, nil
}

func (f *function) detail() string {
	str := fmt.Sprintf("Func[%d]:\n", len(f.types))
	for i := 0; i < len(f.types); i++ {
		str += fmt.Sprintf(" - func[%d] sig=%d\n", i, f.types[i])
	}
	return str
}
