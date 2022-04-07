package instruction

type Call struct {
	Imm uint32
}

func (*Call) Opcode() Opcode {
	return CALL
}

func (c *Call) imm() any {
	return c.Imm
}

func (*Call) String() string {
	return "call"
}

type CallIndirect struct {
	Imm CallIndirectImm
}

type CallIndirectImm struct {
	TypeIndex uint32
	reserved  bool
}

func (*CallIndirect) Opcode() Opcode {
	return CALL_INDIRECT
}

func (ci *CallIndirect) imm() any {
	return ci.Imm
}

func (*CallIndirect) String() string {
	return "call_indirect"
}
