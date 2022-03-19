package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type mod struct {
	version  uint32
	custom   *custom
	typ      *typ
	imports  *imports
	function *function
	table    *table
	memory   *memory
	global   *global
	export   *export
	start    *start
	element  *element
	code     *code
	data     *data
}

func (m *mod) build() (*structure.Module, error) {
	sm := &structure.Module{Version: m.version}
	if m.typ != nil {
		sm.Types = m.typ.entries
	}
	if m.imports != nil {
		sm.Imports = make([]*structure.Import, 0, len(m.imports.entries))
		for _, i := range m.imports.entries {
			desc, err := fromImportEntry(i.kind, i.typ)
			if err != nil {
				return nil, fmt.Errorf("module build: %w", err)
			}
			sm.Imports = append(sm.Imports, &structure.Import{
				Module: string(i.moduleName),
				Name:   string(i.fieldString),
				Desc:   desc,
			})
		}
	}
	if m.function != nil {
		sm.Functions = make([]*structure.Function, 0, len(m.function.types))
		if sm.Imports != nil {
			for _, imp := range sm.Imports {
				if imp.Desc.Type == structure.DescTypeFunc {
					sm.Functions = append(sm.Functions, &structure.Function{Type: imp.Desc.Func, Imported: true})
				}
			}
		}
		for i, typ := range m.function.types {
			f := &structure.Function{Type: typ}
			f.Locals = make([]types.ValueType, 0, len(m.code.bodies[i].locals))
			for _, l := range m.code.bodies[i].locals {
				f.Locals = append(f.Locals, l.typ)
			}
			f.Body = m.code.bodies[i].code
			// sm.Functions[typ] = f
			sm.Functions = append(sm.Functions, f)
		}
		// for _, i := range sm.Imports {
		// 	if i.Desc.Type == structure.DescTypeFunc {
		// 		sm.Functions[i.Desc.Func] = &structure.Function{Type: i.Desc.Func}
		// 	}
		// }
	}
	if m.table != nil {
		sm.Tables = make([]*structure.Table, 0, len(m.table.entries))
		for _, t := range m.table.entries {
			sm.Tables = append(sm.Tables, &structure.Table{Type: t})
		}
	}
	if m.memory != nil {
		sm.Memories = make([]*structure.Memory, 0, len(m.memory.entries))
		for _, mem := range m.memory.entries {
			sm.Memories = append(sm.Memories, &structure.Memory{Type: mem})
		}
	}
	if m.global != nil {
		sm.Globals = make([]*structure.Global, 0, len(m.global.globals))
		for _, g := range m.global.globals {
			instr, err := instruction.Decode(bytes.NewBuffer(g.init))
			if err != nil {
				return nil, fmt.Errorf("module build: global init expr: %w", err)
			}
			sm.Globals = append(sm.Globals, &structure.Global{Type: g.typ, Init: instr})
		}
	}
	if m.element != nil {
		sm.Elements = make([]*structure.Element, 0, len(m.element.entries))
		for _, e := range m.element.entries {
			instr, err := instruction.Decode(bytes.NewBuffer(e.offset))
			if err != nil {
				return nil, fmt.Errorf("module build: element init expr: %w", err)
			}
			sm.Elements = append(sm.Elements, &structure.Element{
				Type:       types.ElemTypeFuncref,
				TableIndex: e.index,
				Offset:     instr,
				Init:       e.elems,
			})
		}
	}
	if m.data != nil {
		sm.Datas = make([]*structure.Data, 0, len(m.data.entries))
		for _, d := range m.data.entries {
			instr, err := instruction.Decode(bytes.NewBuffer(d.offset))
			if err != nil {
				return nil, fmt.Errorf("module build: data init expr: %w", err)
			}
			sm.Datas = append(sm.Datas, &structure.Data{
				Init:        d.data,
				MemoryIndex: d.index,
				Offset:      instr,
			})
		}
	}
	if m.start != nil {
		sm.Start = &structure.Start{Index: m.start.index}
	}
	if m.export != nil {
		sm.Exports = make([]*structure.Export, 0, len(m.export.entries))
		for _, e := range m.export.entries {
			sm.Exports = append(sm.Exports, &structure.Export{
				Name: string(e.fieldString),
				Desc: &structure.ExportDesc{
					Type: structure.DescType(e.kind),
					Val:  e.index,
				},
			})
		}
	}
	return sm, nil
}
