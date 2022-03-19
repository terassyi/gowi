package instance

import (
	"fmt"

	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type Module struct {
	Types      []*types.FuncType
	FuncAddrs  []*Function
	TableAddrs []*Table
	MemAddrs   []*Memory
	GlobalAddr []*Global
	// ElemAddrs  []*Element
	// DataAddrs  []*Data
	Exports []*Export
}

func New(mod *structure.Module) (*Module, error) {
	m := &Module{}
	m.Types = mod.Types
	funcs := newFunctions(mod)
	m.FuncAddrs = funcs
	m.TableAddrs = newTables(mod)
	m.MemAddrs = newMemories(mod)
	for _, e := range mod.Elements {
		table := m.TableAddrs[e.TableIndex]
		offset, err := evaluateConstInstr(e.Offset)
		if err != nil {
			return nil, fmt.Errorf("New module instance: %w", err)
		}
		if err := table.grow(e.Type, GetVal[int32](offset), e.Init, m.FuncAddrs); err != nil {
			return nil, fmt.Errorf("New module instance: %w", err)
		}
	}
	for _, d := range mod.Datas {
		mem := m.MemAddrs[d.MemoryIndex]
		offset, err := evaluateConstInstr(d.Offset)
		if err != nil {
			return nil, fmt.Errorf("New module instance: %w", err)
		}
		if err := mem.grow(GetVal[int32](offset), d.Init); err != nil {
			return nil, fmt.Errorf("New module instance: %w", err)
		}
	}
	exports, err := newExports(mod, m.FuncAddrs, m.TableAddrs, m.MemAddrs, m.GlobalAddr)
	if err != nil {
		return nil, fmt.Errorf("New module instance: %w", err)
	}
	m.Exports = exports
	for _, f := range m.FuncAddrs {
		f.Module = m
	}
	return m, nil
}
