package instruction

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImm_None(t *testing.T) {
	for _, d := range []Instruction{
		&Unreachable{},
		&Nop{},
		&Else{},
		&End{},
		&Return{},
		&Drop{},
		&Select{},
		&I32Add{},
		&I32Sub{},
		&I32Mul{},
	} {
		res := Imm[None](d)
		assert.Equal(t, NoImm, res)
	}
}

func TestImm_Uint32(t *testing.T) {
	for _, d := range []struct {
		instr    Instruction
		expected uint32
	}{
		{instr: &Call{Imm: 0}, expected: uint32(0)},
		{instr: &GetLocal{Imm: 0}, expected: uint32(0)},
		{instr: &GetLocal{Imm: 1}, expected: uint32(1)},
		{instr: &SetLocal{Imm: 1}, expected: uint32(1)},
		{instr: &TeeLocal{Imm: 0xff}, expected: uint32(0xff)},
		{instr: &F32Const{Imm: 0xeb}, expected: uint32(0xeb)},
	} {
		res := Imm[uint32](d.instr)
		assert.Equal(t, d.expected, res)
	}
}

func TestImm_CallIndirectImm(t *testing.T) {
	imm := CallIndirectImm{
		TypeIndex: uint32(0),
		reserved:  true,
	}
	instr := &CallIndirect{Imm: imm}
	res := Imm[CallIndirectImm](instr)
	assert.Equal(t, imm, res)
}

func TestImm_BrTableImm(t *testing.T) {
	imm := BrTableImm{
		TargetTable:   []uint32{0x00, 0x01, 0xff},
		DefaultTarget: uint32(0x00),
	}
	instr := &BrTable{Imm: imm}
	res := Imm[BrTableImm](instr)
	assert.Equal(t, imm, res)
}
