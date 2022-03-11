package structure

import (
	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"
)

type Module struct {
	Version   uint32
	Types     []*types.FuncType
	Functions []*Function
	Tables    []*Table
	Memories  []*Memory
	Global    []*Global
	Elements  []*Element
	Datas     []*Data
	Start     *Start
	Imports   []*Import
	Exports   []*Export
}

type Function struct {
	Type   uint32 // typeidx
	Locals []types.ValueType
	Body   []instruction.Instruction
}

type Table struct {
	Type *types.TableType
}

type Memory struct {
	Type *types.MemoryType
}

type Global struct {
	Type *types.GlobalType
	Init instruction.Instruction
}

// https://webassembly.github.io/spec/core/syntax/modules.html#element-segments
// only support active.
type Element struct {
	Type       types.ElemType
	TableIndex uint32
	Offset     instruction.Instruction
}

// https://webassembly.github.io/spec/core/syntax/modules.html#data-segments
// only support active
type Data struct {
	Init        []byte
	MemoryIndex uint32
	Offset      instruction.Instruction
}

type Start struct {
	Index uint32
}

type Import struct {
	Module string
	Name   string
	Desc   *ImportDesc
}

type ImportDesc struct {
	Type   DescType
	Func   uint32
	Table  *types.TableType
	Mem    *types.MemoryType
	Global *types.GlobalType
}

type Export struct {
	Name string
	Desc *ExportDesc
}

type ExportDesc struct {
	Type DescType
	Val  uint32
}

type DescType uint8

const (
	DescTypeFunc   DescType = 0
	DescTypeTable  DescType = 1
	DescTypeMemory DescType = 2
	DescTypeGlobal DescType = 3
)

// func (m *Module) Dump() string {
// 	str := fmt.Sprintf("WASM file format: %x\n\n", m.Version)
// 	if m.Custom != nil {
// 		str += fmt.Sprintf("%s: not implemented.\n", m.Custom.Code())
// 	}
// 	if m.Type != nil {
// 		str += fmt.Sprintf("%s : count=0x%04x\n", m.Type.Code(), len(m.Type.Entries))
// 	}
// 	if m.Import != nil {
// 		str += fmt.Sprintf("%s : count=0x%04x\n", m.Import.Code(), len(m.Import.Entries))
// 	}
// 	if m.Function != nil {
// 		str += fmt.Sprintf("%s : count=0x%04x\n", m.Function.Code(), len(m.Function.Types))
// 	}
// 	if m.Table != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Table.Code(), len(m.Table.Entries))
// 	}
// 	if m.Memory != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Memory.Code(), len(m.Memory.Entries))
// 	}
// 	if m.Global != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Global.Code(), len(m.Global.Globals))
// 	}
// 	if m.Export != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Export.Code(), len(m.Export.Entries))
// 	}
// 	if m.Start != nil {
// 		str += fmt.Sprintf("%s: index=%d\n", m.Start.Code(), m.Start.Index)
// 	}
// 	if m.Element != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Element.Code(), len(m.Element.Entries))
// 	}
// 	if m.Code != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Code.Code(), len(m.Code.Bodies))
// 	}
// 	if m.Data != nil {
// 		str += fmt.Sprintf("%s: count=0x%04x\n", m.Data.Code(), len(m.Data.Entries))
// 	}
// 	return str
// }
//
// func (m *Module) DumpDetail() (string, error) {
// 	str := fmt.Sprintf("WASM file format: %x\n\n", m.Version)
// 	if m.Custom != nil {
// 		str += m.Custom.Detail()
// 		str += "\n"
// 	}
// 	if m.Type != nil {
// 		str += m.Type.Detail()
// 		str += "\n"
// 	}
// 	if m.Import != nil {
// 		str += m.Import.Detail()
// 		str += "\n"
// 	}
// 	if m.Function != nil {
// 		str += m.Function.Detail()
// 		str += "\n"
// 	}
// 	if m.Table != nil {
// 		str += m.Table.Detail()
// 		str += "\n"
// 	}
// 	if m.Memory != nil {
// 		str += m.Memory.Detail()
// 		str += "\n"
// 	}
// 	if m.Global != nil {
// 		s, err := m.Global.Detail()
// 		if err != nil {
// 			return "", err
// 		}
// 		str += s
// 		str += "\n"
// 	}
// 	if m.Export != nil {
// 		str += m.Export.Detail()
// 		str += "\n"
// 	}
// 	if m.Start != nil {
// 		str += m.Start.Detail()
// 		str += "\n"
// 	}
// 	if m.Element != nil {
// 		s, err := m.Element.Detail()
// 		if err != nil {
// 			return "", err
// 		}
// 		str += s
// 		str += "\n"
// 	}
// 	if m.Code != nil {
// 		str += m.Code.Detail()
// 		str += "\n"
// 	}
// 	if m.Data != nil {
// 		s, err := m.Data.Detail()
// 		if err != nil {
// 			return "", err
// 		}
// 		str += s
// 		str += "\n"
// 	}
// 	return str, nil
// }
//
