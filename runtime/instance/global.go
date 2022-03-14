package instance

import (
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

type Global struct {
	Type  types.ValueType
	Value value.Value
}
