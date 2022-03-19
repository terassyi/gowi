package instruction

import "github.com/terassyi/gowi/types"

type Unreachable struct{}

func (*Unreachable) Opcode() Opcode {
	return UNREACHABLE
}

func (*Unreachable) imm() any {
	return NoImm
}

type Nop struct{}

func (*Nop) Opcode() Opcode {
	return NOP
}

func (*Nop) imm() any {
	return NoImm
}

type Block struct {
	Imm types.BlockType
}

func (*Block) Opcode() Opcode {
	return BLOCK
}

func (b *Block) imm() any {
	return b.Imm
}

type Loop struct {
	Imm types.BlockType
}

func (*Loop) Opcode() Opcode {
	return LOOP
}

func (l *Loop) imm() any {
	return l.Imm
}

type If struct {
	Imm types.BlockType
}

func (*If) Opcode() Opcode {
	return IF
}

func (i *If) imm() any {
	return i.Imm
}

type Else struct{}

func (*Else) Opcode() Opcode {
	return ELSE
}

func (*Else) imm() any {
	return NoImm
}

type End struct{}

func (*End) Opcode() Opcode {
	return END
}

func (*End) imm() any {
	return NoImm
}

type Br struct {
	Imm uint32
}

func (*Br) Opcode() Opcode {
	return BR
}

func (b *Br) imm() any {
	return b.Imm
}

type BrIf struct {
	Imm uint32
}

func (*BrIf) Opcode() Opcode {
	return BR_IF
}

func (bi *BrIf) imm() any {
	return bi.Imm
}

type BrTableImm struct {
	TargetTable   []uint32
	DefaultTarget uint32
}

type BrTable struct {
	Imm BrTableImm
}

func (*BrTable) Opcode() Opcode {
	return BR_TABLE
}

func (bt *BrTable) imm() any {
	return bt.Imm
}

type Return struct{}

func (*Return) Opcode() Opcode {
	return RETURN
}

func (*Return) imm() any {
	return NoImm
}
