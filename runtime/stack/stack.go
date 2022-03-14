package stack

import (
	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/runtime/instance"
	"github.com/terassyi/gowi/runtime/value"
)

// https://webassembly.github.io/spec/core/exec/runtime.html#stack
type Stack struct {
}

type Label struct {
	Instructions []instruction.Instruction
	n            uint8
}

type Frame struct {
	Locals []value.Value
	Module *instance.Module
}
