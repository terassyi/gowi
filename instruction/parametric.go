package instruction

type Drop struct{}

func (*Drop) Opcode() Opcode {
	return DROP
}

func (*Drop) imm() any {
	return NoImm
}

func (*Drop) String() string {
	return "drop"
}

func (*Drop) ImmString() string {
	return ""
}

type Select struct{}

func (*Select) Opcode() Opcode {
	return SELECT
}

func (*Select) imm() any {
	return NoImm
}

func (*Select) String() string {
	return "select"
}

func (*Select) ImmString() string {
	return ""
}
