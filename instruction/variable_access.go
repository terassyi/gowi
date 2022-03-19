package instruction

type GetLocal struct {
	Imm uint32
}

func (*GetLocal) Opcode() Opcode {
	return GET_LOCAL
}

func (gl *GetLocal) imm() any {
	return gl.Imm
}

type SetLocal struct {
	Imm uint32
}

func (*SetLocal) Opcode() Opcode {
	return SET_LOCAL
}

func (sl *SetLocal) imm() any {
	return sl.Imm
}

type TeeLocal struct {
	Imm uint32
}

func (*TeeLocal) Opcode() Opcode {
	return TEE_LOCAL
}

func (tl *TeeLocal) imm() any {
	return tl.Imm
}

type GetGlobal struct {
	Imm uint32
}

func (*GetGlobal) Opcode() Opcode {
	return GET_GLOBAL
}

func (gg *GetGlobal) imm() any {
	return gg.Imm
}

type SetGlobal struct {
	Imm uint32
}

func (*SetGlobal) Opcode() Opcode {
	return SET_GLOBAL
}

func (sg *SetGlobal) imm() any {
	return sg.Imm
}
