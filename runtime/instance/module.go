package instance

import (
	"github.com/terassyi/gowi/types"
)

type Module struct {
	Types      []types.FuncType
	FuncAddrs  []*Function
	TableAddrs []*Table
	MemAddrs   []*Memory
	GlobalAddr []*Global
	ElemAddrs  []*Element
	DataAddrs  []*Data
	Exports    []*Export
}

type ReferenceTypeSet interface {
	*Function
}
