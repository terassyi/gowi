package instruction

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/types"
)

type MemoryImm struct {
	Flags  uint32
	Offset uint32
}

type I32Load struct{ Imm MemoryImm }

func (*I32Load) Opcode() Opcode {
	return I32_LOAD
}

func (i *I32Load) imm() any {
	return i.Imm
}

func (*I32Load) String() string {
	return "i32.load"
}

type I64Load struct{ Imm MemoryImm }

func (*I64Load) Opcode() Opcode {
	return I64_LOAD
}

func (i *I64Load) imm() any {
	return i.Imm
}

func (*I64Load) String() string {
	return "i64.load"
}

type I32Load8S struct{ Imm MemoryImm }

func (*I32Load8S) Opcode() Opcode {
	return I32_LOAD8_S
}

func (i *I32Load8S) imm() any {
	return i.Imm
}

func (*I32Load8S) String() string {
	return "i32.load8_s"
}

type I64Load8S struct{ Imm MemoryImm }

func (*I64Load8S) Opcode() Opcode {
	return I64_LOAD8_S
}

func (i *I64Load8S) imm() any {
	return i.Imm
}

func (*I64Load8S) String() string {
	return "i64.load8_s"
}

type I32Load8U struct{ Imm MemoryImm }

func (*I32Load8U) Opcode() Opcode {
	return I32_LOAD8_U
}

func (i *I32Load8U) imm() any {
	return i.Imm
}

func (*I32Load8U) String() string {
	return "i32.load8_u"
}

type I64Load8U struct{ Imm MemoryImm }

func (*I64Load8U) Opcode() Opcode {
	return I64_LOAD8_U
}

func (i *I64Load8U) imm() any {
	return i.Imm
}

func (*I64Load8U) String() string {
	return "i64.load8_u"
}

type I32Load16S struct{ Imm MemoryImm }

func (*I32Load16S) Opcode() Opcode {
	return I32_LOAD16_S
}

func (i *I32Load16S) imm() any {
	return i.Imm
}

func (*I32Load16S) String() string {
	return "i32.load16_s"
}

type I64Load16S struct{ Imm MemoryImm }

func (*I64Load16S) Opcode() Opcode {
	return I64_LOAD16_S
}

func (i *I64Load16S) imm() any {
	return i.Imm
}

func (*I64Load16S) String() string {
	return "i64.load16_s"
}

type I32Load16U struct{ Imm MemoryImm }

func (*I32Load16U) Opcode() Opcode {
	return I32_LOAD16_U
}

func (i *I32Load16U) imm() any {
	return i.Imm
}

func (*I32Load16U) String() string {
	return "i32.load16_u"
}

type I64Load16U struct{ Imm MemoryImm }

func (*I64Load16U) Opcode() Opcode {
	return I64_LOAD16_U
}

func (i *I64Load16U) imm() any {
	return i.Imm
}

func (*I64Load16U) String() string {
	return "i64.load16_u"
}

type I64Load32S struct{ Imm MemoryImm }

func (*I64Load32S) Opcode() Opcode {
	return I64_LOAD32_S
}

func (i *I64Load32S) imm() any {
	return i.Imm
}

func (*I64Load32S) String() string {
	return "i64.load32_s"
}

type I64Load32U struct{ Imm MemoryImm }

func (*I64Load32U) Opcode() Opcode {
	return I64_LOAD32_U
}

func (i *I64Load32U) imm() any {
	return i.Imm
}

func (*I64Load32U) String() string {
	return "i64.load32_u"
}

type I32Store struct{ Imm MemoryImm }

func (*I32Store) Opcode() Opcode {
	return I32_STORE
}

func (i *I32Store) imm() any {
	return i.Imm
}

func (*I32Store) String() string {
	return "i32.store"
}

type I64Store struct{ Imm MemoryImm }

func (*I64Store) Opcode() Opcode {
	return I64_STORE
}

func (i *I64Store) imm() any {
	return i.Imm
}

func (*I64Store) String() string {
	return "i64.store"
}

type I32Store8 struct{ Imm MemoryImm }

func (*I32Store8) Opcode() Opcode {
	return I32_STORE8
}

func (i *I32Store8) imm() any {
	return i.Imm
}

func (*I32Store8) String() string {
	return "i32.store8"
}

type I64Store8 struct{ Imm MemoryImm }

func (*I64Store8) Opcode() Opcode {
	return I64_STORE8
}

func (i *I64Store8) imm() any {
	return i.Imm
}

func (*I64Store8) String() string {
	return "i64.store8"
}

type I32Store16 struct{ Imm MemoryImm }

func (*I32Store16) Opcode() Opcode {
	return I32_STORE16
}

func (i *I32Store16) imm() any {
	return i.Imm
}

func (*I32Store16) String() string {
	return "i32.store16"
}

type I64Store16 struct{ Imm MemoryImm }

func (*I64Store16) Opcode() Opcode {
	return I64_STORE16
}

func (i *I64Store16) imm() any {
	return i.Imm
}

func (*I64Store16) String() string {
	return "i64.store16"
}

type I64Store32 struct{ Imm MemoryImm }

func (*I64Store32) Opcode() Opcode {
	return I64_STORE32
}

func (i *I64Store32) imm() any {
	return i.Imm
}

func (*I64Store32) String() string {
	return "i64.store32"
}

func newMemImm(buf *bytes.Buffer) (*MemoryImm, error) {
	flags, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("mem imm: %w", err)
	}
	offset, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("mem imm: %w", err)
	}
	return &MemoryImm{Flags: uint32(flags), Offset: uint32(offset)}, nil
}
