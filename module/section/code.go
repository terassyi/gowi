package section

import (
	"bytes"
	"fmt"

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
	Code   []byte
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
	return &FunctionBody{
		Locals: locals,
		Code:   code[:len(code)-1],
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
