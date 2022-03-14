package instance

import (
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/types"
)

type Table struct {
	Type types.TableType
	Elem []value.Reference
}
