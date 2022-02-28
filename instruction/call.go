package instruction

type Call struct {
	Imm uint32
}

func (*Call) Opcode() Opcode {
	return CALL
}

type CallIndirect struct {
	Imm *CallIndirectImm
}

type CallIndirectImm struct {
	TypeIndex uint32
	reserved  bool
}

func (*CallIndirect) Opcode() Opcode {
	return CALL_INDIRECT
}
