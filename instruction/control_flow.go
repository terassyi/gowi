package instruction

import "github.com/terassyi/gowi/types"

type Unreachable struct{}

func (*Unreachable) Opcode() Opcode {
	return UNREACHABLE
}

type Nop struct{}

func (*Nop) Opcode() Opcode {
	return NOP
}

type Block struct {
	Imm types.BlockType
}

func (*Block) Opcode() Opcode {
	return BLOCK
}

type Loop struct {
	Imm types.BlockType
}

func (*Loop) Opcode() Opcode {
	return LOOP
}

type If struct {
	Imm types.BlockType
}

func (*If) Opcode() Opcode {
	return IF
}

type Else struct{}

func (*Else) Opcode() Opcode {
	return ELSE
}

type End struct{}

func (*End) Opcode() Opcode {
	return END
}

type Br struct {
	Imm uint32
}

func (*Br) Opcode() Opcode {
	return BR
}

type BrIf struct {
	Imm uint32
}

func (*BrIf) Opcode() Opcode {
	return BR_IF
}

type BrTableImm struct {
	TargetTable   []uint32
	DefaultTarget uint32
}

type BrTable struct {
	Imm *BrTableImm
}

func (*BrTable) Opcode() Opcode {
	return BR_TABLE
}

type Return struct{}

func (*Return) Opcode() Opcode {
	return RETURN
}
