package decoder

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"
)

type code struct {
	// count  uint32
	bodies []*functionBody
}

type functionBody struct {
	// BodySize   uint32
	// LocalCount uint32
	locals []*localEntry
	code   []instruction.Instruction
}

type localEntry struct {
	count uint32
	typ   types.ValueType
}

func newCode(payload []byte) (*code, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewCode: decode count: %w", err)
	}
	funcBodys := make([]*functionBody, 0, int(count))
	for i := 0; i < int(count); i++ {
		f, err := newFunctionBody(buf)
		if err != nil {
			return nil, fmt.Errorf("NewCode: decode function_body: %w", err)
		}
		funcBodys = append(funcBodys, f)
	}
	return &code{
		bodies: funcBodys,
	}, nil
}

func newFunctionBody(buf *bytes.Buffer) (*functionBody, error) {
	_, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("newFunctionBody: decode body_size: %w", err)
	}
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("newFunctionBody: decode local_count: %w", err)
	}
	locals := make([]*localEntry, 0, int(count))
	for i := 0; i < int(count); i++ {
		l, err := newLocalEntry(buf)
		if err != nil {
			return nil, fmt.Errorf("newFunctionBody: decode locals: %w", err)
		}
		locals = append(locals, l)
	}
	code, err := buf.ReadBytes(END)
	if err != nil {
		return nil, fmt.Errorf("newFunctionBody: decode code: %w", err)
	}
	codeBuf := bytes.NewBuffer(code)
	codes := make([]instruction.Instruction, 0, 1024)
	for codeBuf.Len() > 0 {
		c, err := instruction.Decode(codeBuf)
		if err != nil {
			return nil, fmt.Errorf("newFunctionBody: decode instruction: %w", err)
		}
		codes = append(codes, c)
	}
	return &functionBody{
		locals: locals,
		code:   codes,
	}, nil
}

func newLocalEntry(buf *bytes.Buffer) (*localEntry, error) {
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("newLocalEntry: decode count: %w", err)
	}
	t, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("newLocalEntry: decode type: %w", err)
	}
	typ, err := types.NewValueType(t)
	if err != nil {
		return nil, fmt.Errorf("newLocalEntry: decode type: %w", err)
	}
	return &localEntry{
		count: uint32(count),
		typ:   typ,
	}, nil
}

func (c *code) detail() string {
	str := fmt.Sprintf("Code[%d]:\n", len(c.bodies))
	for i := 0; i < len(c.bodies); i++ {
		str += fmt.Sprintf(" - func[%d] instruction size=%d\n", i, len(c.bodies[i].code))
	}
	return str
}
