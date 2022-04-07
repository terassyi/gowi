package instance

import (
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

type Element struct {
	Type types.ReferenceType
	Elem []value.Reference // ref
}
