package instruction

import "fmt"

type GetLocal struct {
	Imm uint32
}

func (*GetLocal) Opcode() Opcode {
	return GET_LOCAL
}

func (gl *GetLocal) imm() any {
	return gl.Imm
}

func (*GetLocal) String() string {
	return "get_local"
}

func (g *GetLocal) ImmString() string {
	return fmt.Sprintf("$%d", g.Imm)
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

func (*SetLocal) String() string {
	return "set_local"
}

func (s *SetLocal) ImmString() string {
	return fmt.Sprintf("$%d", s.Imm)
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

func (*TeeLocal) String() string {
	return "tee_local"
}

func (t *TeeLocal) ImmString() string {
	return fmt.Sprintf("$%d", t.Imm)
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

func (*GetGlobal) String() string {
	return "get_global"
}

func (g *GetGlobal) ImmString() string {
	return fmt.Sprintf("$%d", g.Imm)
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

func (*SetGlobal) String() string {
	return "set_global"
}

func (s *SetGlobal) ImmString() string {
	return fmt.Sprintf("$%d", s.Imm)
}
