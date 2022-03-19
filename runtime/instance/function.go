package instance

import (
	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

// https://webassembly.github.io/spec/core/exec/runtime.html#function-instances
// TODO hostfunc
type Function struct {
	Type   *types.FuncType
	Module *Module
	Code   *structure.Function
}

func newFunctions(mod *structure.Module) []*Function {
	funcs := make([]*Function, 0, len(mod.Functions))
	for _, f := range mod.Functions {
		funcs = append(funcs, &Function{
			Type: mod.Types[f.Type],
			// Moudle: after instanciating function instances, creates references to a module instance.
			Code: f,
		})
	}
	return funcs
}

func (*Function) RefType() value.ReferenceType {
	return value.RefTypeFunc
}

func (*Function) ExternalValueType() ExternValueType {
	return ExternValTypeFunc
}
