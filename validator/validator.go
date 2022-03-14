package validator

import (
	"errors"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/structure"
	"github.com/terassyi/gowi/types"
)

var (
	ContentTypeIsNotMatched  error = errors.New("content type is not matched")
	NotEmptyFuncBodyExpected error = errors.New("empty function body is not expected")
	TooManyIndexSpace        error = errors.New("too many index space")
)

type Validator struct {
	mod *structure.Module
	ctx *context
}

func New(mod *structure.Module) (*Validator, error) {
	ctx, err := newContext(mod)
	if err != nil {
		return nil, fmt.Errorf("Validator new: %w", err)
	}
	return &Validator{
		mod: mod,
		ctx: ctx,
	}, nil
}

func (v *Validator) Validate() (bool, error) {
	if v.mod.Tables != nil {
		if len(v.ctx.tables) > 1 {
			return false, fmt.Errorf("Validate error: %w", TooManyIndexSpace)
		}
		for _, t := range v.ctx.tables {
			if err := validateTable(t); err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
		}
	}
	if v.mod.Memories != nil {
		if len(v.ctx.memories) > 1 {
			return false, fmt.Errorf("Validate error: %w", TooManyIndexSpace)
		}
		for _, m := range v.ctx.memories {
			if err := validateMemory(m); err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
		}
	}
	if v.mod.Globals != nil {
		for _, g := range v.mod.Globals {
			if err := validateGlobal(g); err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
		}
	}
	if v.mod.Functions != nil {
		for _, f := range v.mod.Functions {
			if err := validateFunction(f); err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
		}
	}
	if v.mod.Start != nil {
		startFunc := v.ctx.functions[v.mod.Start.Index]
		if !startFunc.Params.IsEmpty() || !startFunc.Returns.IsEmpty() {
			return false, fmt.Errorf("Validate error: the start function is expected to have type [] -> []")
		}
	}
	if v.mod.Exports != nil {
		dup := map[string]bool{}
		for _, e := range v.mod.Exports {
			// export name duplication check
			if _, ok := dup[e.Name]; ok {
				return false, fmt.Errorf("Validate error: duplicate export %s", e.Name)
			}
			dup[e.Name] = true
			// check export type
			switch e.Desc.Type {
			case structure.DescTypeFunc:
				if _, err := v.ctx.reqiureFunc(e.Desc.Val); err != nil {
					return false, fmt.Errorf("Validate error: export desc: %w", err)
				}
			case structure.DescTypeTable:
				if _, err := v.ctx.requireTable(e.Desc.Val); err != nil {
					return false, fmt.Errorf("Validate error: export desc: %w", err)
				}
			case structure.DescTypeMemory:
				if _, err := v.ctx.requireMemory(e.Desc.Val); err != nil {
					return false, fmt.Errorf("Validate error: export desc: %w", err)
				}
			case structure.DescTypeGlobal:
				if _, err := v.ctx.requireGlobal(e.Desc.Val); err != nil {
					return false, fmt.Errorf("Validate error: export desc: %w", err)
				}
			default:
				return false, fmt.Errorf("Validate error: invalid export type")

			}
		}
	}
	if v.mod.Imports != nil {
		for _, i := range v.mod.Imports {
			switch i.Desc.Type {
			case structure.DescTypeFunc:
				if _, err := v.ctx.reqiureFunc(i.Desc.Func); err != nil {
					return false, fmt.Errorf("Validate error: import desc: %w", err)
				}
			case structure.DescTypeTable:
				if err := validateTable(i.Desc.Table); err != nil {
					return false, fmt.Errorf("Validate error: import desc: %w", err)
				}
			case structure.DescTypeMemory:
				if err := validateMemory(i.Desc.Mem); err != nil {
					return false, fmt.Errorf("Validate error: import desc: %w", err)
				}
			case structure.DescTypeGlobal:
				if i.Desc.Global == nil {
					return false, fmt.Errorf("Validate error: import desc: global param is nil")
				}
			default:
				return false, fmt.Errorf("Validate error: invalid import type")
			}
		}
	}
	if v.mod.Datas != nil {
		for _, d := range v.mod.Datas {
			if _, err := v.ctx.requireMemory(d.MemoryIndex); err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
			typ, err := instruction.GetConstType(d.Offset)
			if err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
			if typ != types.I32 {
				return false, fmt.Errorf("Validate error: global init instruction is not i32.const: given %s", typ)
			}
		}
	}
	if v.mod.Elements != nil {
		for _, e := range v.mod.Elements {
			if _, err := v.ctx.requireTable(e.TableIndex); err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
			typ, err := instruction.GetConstType(e.Offset)
			if err != nil {
				return false, fmt.Errorf("Validate error: %w", err)
			}
			if typ != types.I32 {
				return false, fmt.Errorf("Validate error: global init instruction is not i32.const: given %s", typ)
			}
		}
	}

	return true, nil
}

func validateTable(tableType *types.TableType) error {
	return tableType.Limits.Validate()
}

func validateMemory(memType *types.MemoryType) error {
	return memType.Limits.Validate()
}

func validateGlobal(global *structure.Global) error {
	t, err := instruction.GetConstType(global.Init)
	if err != nil {
		return fmt.Errorf("validateGlobal: %w", err)
	}
	if t != global.Type.ContentType {
		return fmt.Errorf("validateGlobal: %w: want=%x", ContentTypeIsNotMatched, t)
	}
	return nil
}

func validateFunction(f *structure.Function) error {
	if f.Imported {
		return nil
	}
	if f.Body == nil || len(f.Body) == 0 {
		return fmt.Errorf("validateFunction: %w", NotEmptyFuncBodyExpected)
	}
	v := &funcValidator{}
	for i, instr := range f.Body {
		if err := v.step(instr); err != nil {
			return fmt.Errorf("validateaFunction: validation error at %d th instruction in function(%d)", i, f.Type)
		}
	}
	return nil
}
