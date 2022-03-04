package module

import (
	"fmt"

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

func (m *Module) Dump() string {
	str := fmt.Sprintf("WASM file format: %x\n\n", m.Version)
	if m.Custom != nil {
		str += fmt.Sprintf("%s: not implemented.\n", m.Custom.Code())
	}
	if m.Type != nil {
		str += fmt.Sprintf("%s : count=0x%04x\n", m.Type.Code(), len(m.Type.Entries))
	}
	if m.Import != nil {
		str += fmt.Sprintf("%s : count=0x%04x\n", m.Import.Code(), len(m.Import.Entries))
	}
	if m.Function != nil {
		str += fmt.Sprintf("%s : count=0x%04x\n", m.Function.Code(), len(m.Function.Types))
	}
	if m.Table != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Table.Code(), len(m.Table.Entries))
	}
	if m.Memory != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Memory.Code(), len(m.Memory.Entries))
	}
	if m.Global != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Global.Code(), len(m.Global.Globals))
	}
	if m.Export != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Export.Code(), len(m.Export.Entries))
	}
	if m.Start != nil {
		str += fmt.Sprintf("%s: index=%d\n", m.Start.Code(), m.Start.Index)
	}
	if m.Element != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Element.Code(), len(m.Element.Entries))
	}
	if m.Code != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Code.Code(), len(m.Code.Bodies))
	}
	if m.Data != nil {
		str += fmt.Sprintf("%s: count=0x%04x\n", m.Data.Code(), len(m.Data.Entries))
	}
	return str
}

func (m *Module) DumpDetail() (string, error) {
	str := fmt.Sprintf("WASM file format: %x\n\n", m.Version)
	if m.Custom != nil {
		str += m.Custom.Detail()
		str += "\n"
	}
	if m.Type != nil {
		str += m.Type.Detail()
		str += "\n"
	}
	if m.Import != nil {
		str += m.Import.Detail()
		str += "\n"
	}
	if m.Function != nil {
		str += m.Function.Detail()
		str += "\n"
	}
	if m.Table != nil {
		str += m.Table.Detail()
		str += "\n"
	}
	if m.Memory != nil {
		str += m.Memory.Detail()
		str += "\n"
	}
	if m.Global != nil {
		s, err := m.Global.Detail()
		if err != nil {
			return "", err
		}
		str += s
		str += "\n"
	}
	if m.Export != nil {
		str += m.Export.Detail()
		str += "\n"
	}
	if m.Start != nil {
		str += m.Start.Detail()
		str += "\n"
	}
	if m.Element != nil {
		s, err := m.Element.Detail()
		if err != nil {
			return "", err
		}
		str += s
		str += "\n"
	}
	if m.Code != nil {
		str += m.Code.Detail()
		str += "\n"
	}
	if m.Data != nil {
		s, err := m.Data.Detail()
		if err != nil {
			return "", err
		}
		str += s
		str += "\n"
	}
	return str, nil
}
