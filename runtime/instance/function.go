package instance

import (
	"github.com/terassyi/gowi/types"
)

// https://webassembly.github.io/spec/core/exec/runtime.html#function-instances
// TODO hostfunc
type Function struct {
	Type   types.FuncType
	Module *Module
	// Code   *section.Function
}
