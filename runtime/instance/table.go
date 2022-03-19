package instance

import (
	"fmt"

	"github.com/terassyi/gowi/runtime/value"
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

type Table struct {
	Type  *types.TableType
	Elems []value.Reference
}

func newTables(mod *structure.Module) []*Table {
	tables := make([]*Table, 0, len(mod.Tables))
	for _, t := range mod.Tables {
		tables = append(tables, &Table{
			Type: t.Type,
			// Elems: make([]value.Reference, 0, t.Type.Limits.Min),
			Elems: make([]value.Reference, t.Type.Limits.Min),
		})
	}
	return tables
}

func (*Table) ExternalValueType() ExternValueType {
	return ExternValTypeTable
}

// https://webassembly.github.io/spec/core/exec/modules.html#growing-tables
func (t *Table) grow(typ types.ElemType, offset int32, elems []uint32, funcs []*Function) error {
	if t.Type.ElementType != typ {
		return fmt.Errorf("grow table: type doesn't match")
	}
	if int(offset)+len(elems) > len(t.Elems) {
		if t.Type.Limits.Max != 0 && int(offset)+len(elems) > int(t.Type.Limits.Max) {
			return fmt.Errorf("grow table: exceed the max")
		}
		growDiff := 0
		// if limit max is set
		if t.Type.Limits.Max != 0 {
			growDiff = int(t.Type.Limits.Max) - len(t.Elems)
		} else {
			growDiff = int(offset) + len(elems) - len(t.Elems)
		}
		t.Elems = append(t.Elems, make([]value.Reference, growDiff)...)
	}
	for i, elem := range elems {
		t.Elems[int(offset)+i] = funcs[elem]
	}
	return nil
}

func (t *Table) Len() int {
	return len(t.Elems)
}
