package validator

import (
	"fmt"

	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type context struct {
	types      []*types.FuncType
	functions  []*types.FuncType
	tables     []*types.TableType
	memories   []*types.MemoryType
	globals    []*types.GlobalType
	elements   []types.ReferenceType
	datas      []bool
	locals     []types.ValueType
	labels     []types.ResultType
	returns    []types.ResultType
	references []uint32 // funcidx*
}

func newContext(mod *structure.Module) (*context, error) {
	ctx := &context{}
	if mod.Types != nil {
		ctx.types = mod.Types
		ctx.returns = make([]types.ResultType, 0)
		for _, t := range mod.Types {
			ctx.returns = append(ctx.returns, t.Returns)
		}
	}
	if mod.Functions != nil {
		ctx.functions = make([]*types.FuncType, 0, len(mod.Functions))
		ctx.locals = make([]types.ValueType, 0)
		for _, idx := range mod.Functions {
			fmt.Println(mod.Types)
			fmt.Println(idx)
			ctx.functions = append(ctx.functions, mod.Types[idx.Type])
			ctx.locals = append(ctx.locals, idx.Locals...)
		}
	}
	if mod.Tables != nil {
		ctx.tables = make([]*types.TableType, 0, len(mod.Tables))
		for _, t := range mod.Tables {
			ctx.tables = append(ctx.tables, t.Type)
		}
	}
	if mod.Memories != nil {
		ctx.memories = make([]*types.MemoryType, 0, len(mod.Memories))
		for _, m := range mod.Memories {
			ctx.memories = append(ctx.memories, m.Type)
		}
	}
	if mod.Globals != nil {
		ctx.globals = []*types.GlobalType{}
		for _, g := range mod.Globals {
			ctx.globals = append(ctx.globals, g.Type)
		}
	}
	if mod.Elements != nil {
		ctx.elements = []types.ReferenceType{}
		for _, e := range mod.Elements {
			ctx.elements = append(ctx.elements, types.ReferenceType(e.Type))
		}
	}
	if mod.Datas != nil {
		ctx.datas = []bool{}
		for range mod.Datas {
			ctx.datas = append(ctx.datas, true)
		}
	}
	return ctx, nil
}

func (c *context) reqiureFunc(index uint32) (*types.FuncType, error) {
	if len(c.functions) == 0 || c.functions == nil {
		return nil, fmt.Errorf("function section is not exist.")
	}
	if int(index) > len(c.functions) {
		return nil, fmt.Errorf("function section index is not valid: %d", index)
	}
	fmt.Println(index)
	return c.functions[index], nil
}

func (c *context) requireTable(index uint32) (*types.TableType, error) {
	if len(c.tables) == 0 || c.tables == nil {
		return nil, fmt.Errorf("talbe section is not exist.")
	}
	if int(index) > len(c.tables) {
		return nil, fmt.Errorf("type section index is not valid: %d", index)
	}
	return c.tables[index], nil
}

func (c *context) requireMemory(index uint32) (*types.MemoryType, error) {
	if len(c.memories) == 0 || c.memories == nil {
		return nil, fmt.Errorf("talbe section is not exist.")
	}
	if int(index) > len(c.memories) {
		return nil, fmt.Errorf("type section index is not valid: %d", index)
	}
	return c.memories[index], nil
}

func (c *context) requireGlobal(index uint32) (*types.GlobalType, error) {
	if len(c.globals) == 0 || c.globals == nil {
		return nil, fmt.Errorf("talbe section is not exist.")
	}
	if int(index) > len(c.globals) {
		return nil, fmt.Errorf("type section index is not valid: %d", index)
	}
	return c.globals[index], nil
}
