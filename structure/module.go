package structure

import (
	"errors"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"
)

type Module struct {
	Version   uint32
	Types     []*types.FuncType
	Functions []*Function
	Tables    []*Table
	Memories  []*Memory
	Globals   []*Global
	Elements  []*Element
	Datas     []*Data
	Start     *Start
	Imports   []*Import
	Exports   []*Export
}

type Function struct {
	Type     uint32 // typeidx
	Locals   []types.ValueType
	Body     []instruction.Instruction
	Imported bool
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
	Init       []uint32 // function index
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

var InvalidDesType error = errors.New("Invalid desc type")
