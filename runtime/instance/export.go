package instance

import (
	"fmt"

	"github.com/terassyi/gowi/structure"
)

type Export struct {
	Name  string
	Value ExternalValue
}

func newExports(mod *structure.Module, funcs []*Function, tables []*Table, memories []*Memory, globals []*Global) ([]*Export, error) {
	exports := make([]*Export, 0, len(mod.Exports))
	for _, e := range mod.Exports {
		var val ExternalValue
		switch e.Desc.Type {
		case structure.DescTypeFunc:
			val = funcs[e.Desc.Val]
		case structure.DescTypeTable:
			val = tables[e.Desc.Val]
		case structure.DescTypeMemory:
			val = memories[e.Desc.Val]
		case structure.DescTypeGlobal:
			val = globals[e.Desc.Val]
		default:
			return nil, fmt.Errorf("new export instance: %w", structure.InvalidDesType)
		}
		exports = append(exports, &Export{
			Name:  e.Name,
			Value: val,
		})
	}
	return exports, nil
}
