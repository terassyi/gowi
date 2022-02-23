package module

import (
	"github.com/terassyi/gowi/module/section"
)

const (
	MAJIC_NUMBER uint32 = 0x0061736d // \0asm
	WASM_VERSION uint32 = 0x1
)

type Module struct {
	Version  uint32
	Custom   *section.Custom
	Type     *section.Type
	Import   *section.Import
	Function *section.Function
	Table    *section.Table
	Memory   *section.Memory
	Global   *section.Global
	Export   *section.Export
	Start    *section.Start
	Element  *section.Element
	Code     *section.Code
	Data     *section.Data
}
