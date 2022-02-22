package types

import (
	"bytes"
	"errors"
	"fmt"
)

var (
	InvalidValueType    error = errors.New("Invalid value type")
	InvalidFuncType     error = errors.New("Invalid func type")
	InvalidExternalKind error = errors.New("Invalid external kind value")
	NotImplemented      error = errors.New("Not implemented")
)

type ValueType uint8

const (
	I32     ValueType = 0x7f
	I64     ValueType = 0x7e
	F32     ValueType = 0x7d
	F64     ValueType = 0x7c
	ANYFUNC ValueType = 0x70
	FUNC    ValueType = 0x60
	EMPTY   ValueType = 0x40
)

func NewValueType(v uint8) (ValueType, error) {
	switch v {
	case 0x7f:
		return I32, nil
	case 0x7e:
		return I64, nil
	case 0x7d:
		return F32, nil
	case 0x7c:
		return F64, nil
	case 0x70:
		return ANYFUNC, nil
	case 0x60:
		return FUNC, nil
	case 0x40:
		return EMPTY, nil
	default:
		return ValueType(0x00), InvalidValueType
	}
}

func (v ValueType) String() string {
	switch v {
	case I32:
		return "i32"
	case I64:
		return "i64"
	case F32:
		return "f32"
	case F64:
		return "f64"
	case ANYFUNC:
		return "anyfunc"
	case FUNC:
		return "func"
	case EMPTY:
		return "empty"
	default:
		return "unknown"
	}
}

type BlockType ValueType

type ElemType ValueType // now only allowed anyfunc

type FuncType struct {
	Params  []ValueType
	Returns []ValueType
}

func DecodeFuncType(payload []byte) (*FuncType, int, error) {
	buf := bytes.NewBuffer(payload)
	p, err := buf.ReadByte()
	if err != nil {
		return nil, 0, fmt.Errorf("DecodeFuncType: %w", err)
	}
	params := make([]ValueType, 0, int(p))
	for i := 0; i < int(p); i++ {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, 0, fmt.Errorf("DecodeFuncType: decode params: %w", err)
		}
		param, err := NewValueType(b)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %x", InvalidValueType, payload[i])
		}
		params = append(params, param)
	}
	r, err := buf.ReadByte()
	if err != nil {
		return nil, 0, fmt.Errorf("DecodeFuncType: decode number of returns: %w", err)
	}
	rets := make([]ValueType, 0, int(r))
	for i := 0; i < int(r); i++ {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, 0, fmt.Errorf("DecodeFuncType: decode returns: %w", err)
		}
		ret, err := NewValueType(b)
		if err != nil {
			return nil, 0, fmt.Errorf("%w: %x", InvalidValueType, payload[i])
		}
		rets = append(rets, ret)
	}
	return &FuncType{
		Params:  params,
		Returns: rets,
	}, len(payload) - buf.Len(), nil
}

type GlobalType struct {
	ContentType ValueType
	Mut         bool
}

type TableType struct {
	ElementType ElemType
	Limits      *ResizableLimits
}

func NewTableType(buf *bytes.Buffer) (*TableType, error) {
	elm, _, err := DecodeVarUint32(buf)
	if err != nil {
		return nil, fmt.Errorf("NewTableType: decode elem: %w", err)
	}
	if elm != VarUint32(ANYFUNC) {
		return nil, fmt.Errorf("%w: only allowed anyfunc: %x", NotImplemented, elm)
	}
	r, err := NewResizableLimits(buf)
	if err != nil {
		return nil, fmt.Errorf("NewTableType: decoder resizable_limits: %w", err)
	}
	return &TableType{
		ElementType: ElemType(elm),
		Limits:      r,
	}, nil
}

type MemoryType struct {
	Limits *ResizableLimits
}

func NewMemoryType(buf *bytes.Buffer) (*MemoryType, error) {
	r, err := NewResizableLimits(buf)
	if err != nil {
		return nil, fmt.Errorf("NewMemoryType: decoder resizable_limits: %w", err)
	}
	return &MemoryType{Limits: r}, nil
}

type ExternalKind uint8

const (
	EXTERNAL_KIND_FUNCTION ExternalKind = iota
	EXTERNAL_KIND_TABLE    ExternalKind = iota
	EXTERNAL_KIND_MEMORY   ExternalKind = iota
	EXTERNAL_KIND_GLOBAL   ExternalKind = iota
)

func NewExternalKind(val uint8) (ExternalKind, error) {
	switch val {
	case 0:
		return EXTERNAL_KIND_FUNCTION, nil
	case 1:
		return EXTERNAL_KIND_TABLE, nil
	case 2:
		return EXTERNAL_KIND_MEMORY, nil
	case 3:
		return EXTERNAL_KIND_GLOBAL, nil
	default:
		return 0xff, fmt.Errorf("%w: %x", InvalidExternalKind, val)
	}
}

type ResizableLimits struct {
	Flag    bool
	Initial uint32
	Max     uint32 // if flag is 1
}

func NewResizableLimits(buf *bytes.Buffer) (*ResizableLimits, error) {
	limits := &ResizableLimits{}
	b, err := buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("NewResizableLimits: decode flag: %w", err)
	}
	if b == 1 {
		limits.Flag = true
	}
	b, err = buf.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("NewResizableLimits: decode init: %w", err)
	}
	limits.Initial = uint32(b)
	if limits.Flag {
		b, err := buf.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("NewResizableLimits: decode max: %w", err)
		}
		limits.Max = uint32(b)
	}
	return limits, nil
}
