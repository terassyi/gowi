package instance

import (
	"fmt"

	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type Global struct {
	Type  types.ValueType
	Value value.Value
}

func newGlobals(mod *structure.Module) ([]*Global, error) {
	globals := make([]*Global, 0, len(mod.Globals))
	for _, g := range mod.Globals {
		val, err := evaluateConstInstr(g.Init)
		if err != nil {
			return nil, fmt.Errorf("newGlobal: %w", err)
		}
		globals = append(globals, &Global{
			Type:  g.Type.ContentType,
			Value: val,
		})
	}
	return globals, nil
}

func (*Global) ExternalValueType() ExternValueType {
	return ExternValTypeGlobal
}
