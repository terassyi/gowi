package validator

import "github.com/terassyi/gowi/instruction"

type funcValidator struct {
}

func (v *funcValidator) step(instr instruction.Instruction) error {
	return nil
}
