package instruction

type Drop struct{}

func (*Drop) Opcode() Opcode {
	return DROP
}

type Select struct{}

func (*Select) Opcode() Opcode {
	return SELECT
}
