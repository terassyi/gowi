package section

import (
	"bytes"
	"fmt"

	"github.com/terassyi/gowi/instruction"
	"github.com/terassyi/gowi/types"
)

type Code struct {
	// count  uint32
	Bodies []*FunctionBody
}

type FunctionBody struct {
	// BodySize   uint32
	// LocalCount uint32
	Locals []*LocalEntry
	Code   []instruction.Instruction
}

type LocalEntry struct {
	Count uint32
	Type  types.ValueType
}

func NewCode(payload []byte) (*Code, error) {
	buf := bytes.NewBuffer(payload)
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewCode: decode count: %w", err)
	}
	funcBodys := make([]*FunctionBody, 0, int(count))
	for i := 0; i < int(count); i++ {
		f, err := newFunctionBody(buf)
		if err != nil {
			return nil, fmt.Errorf("NewCode: decode function_body: %w", err)
		}
		funcBodys = append(funcBodys, f)
	}
	return &Code{
		Bodies: funcBodys,
	}, nil
}

func newFunctionBody(buf *bytes.Buffer) (*FunctionBody, error) {
	_, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("newFunctionBody: decode body_size: %w", err)
	}
	count, _, err := types.DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("newFunctionBody: decode local_count: %w", err)
	}
	locals := make([]*LocalEntry, 0, int(count))
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
	return &FunctionBody{
		Locals: locals,
		Code:   codes,
	}, nil
}

func newLocalEntry(buf *bytes.Buffer) (*LocalEntry, error) {
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
	return &LocalEntry{
		Count: uint32(count),
		Type:  typ,
	}, nil
}

func (*Code) Code() SectionCode {
	return CODE
}

func (c *Code) Detail() string {
	str := fmt.Sprintf("%s[%d]:\n", c.Code(), len(c.Bodies))
	for i := 0; i < len(c.Bodies); i++ {
		str += fmt.Sprintf(" - func[%d] instruction size=%d\n", i, len(c.Bodies[i].Code))
	}
	return str
}
