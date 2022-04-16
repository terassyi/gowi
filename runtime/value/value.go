package value

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"unsafe"

	"github.com/terassyi/gowi/types"
)

type Number interface {
	NumType() NumberType
	ValidateValueType(t types.ValueType) bool
	ToValue() Value
}

type NumberTypeSet interface {
	~uint32 | ~uint64 | ~int32 | ~int64 | ~float32 | ~float64
}

type Reference interface {
	RefType() ReferenceType
}

type Value interface {
	ValType() ValueType
	// ExpectNumber() (NumberType, error)
}

type NumberType uint8

const (
	NumTypeI32 NumberType = 0
	NumTypeI64 NumberType = 1
	NumTypeF32 NumberType = 2
	NumTypeF64 NumberType = 3
)

type VectorType uint8

const (
	VecTypeV128 VectorType = 0
)

type ReferenceType uint8

const (
	RefTypeFunc   ReferenceType = 0
	RefTypeExtern ReferenceType = 1
)

type ValueType uint8

const (
	ValTypeNum ValueType = 0
	ValTypeVec ValueType = 1
	ValTypeRef ValueType = 2
)

type I32 uint32

func NewI32[T int32 | uint32](val T) I32 {
	var v uint32
	buff := bytes.NewBuffer(make([]byte, 0, 4))
	binary.Write(buff, binary.BigEndian, val)
	binary.Read(buff, binary.BigEndian, &v)
	return I32(v)
}

func (I32) NumType() NumberType {
	return NumTypeI32
}

func (I32) ValType() ValueType {
	return ValTypeNum
}

func (i I32) ToValue() Value {
	return i
}

func (I32) ValidateValueType(v types.ValueType) bool {
	if v == types.I32 {
		return true
	}
	return false
}

func (i I32) Signed() int32 {
	var v int32
	buff := bytes.NewBuffer(make([]byte, 0, 4))
	binary.Write(buff, binary.BigEndian, uint32(i))
	binary.Read(buff, binary.BigEndian, &v)
	return v
}

func (i I32) Unsigned() uint32 {
	return uint32(i)
}

type I64 uint64

func NewI64[T int64 | uint64](val T) I64 {
	var v uint64
	buff := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(buff, binary.BigEndian, val)
	binary.Read(buff, binary.BigEndian, &v)
	return I64(v)
}

func (i I64) Unsigned() uint64 {
	return uint64(i)
}

func (i I64) Signed() int64 {
	var v int64
	buff := bytes.NewBuffer(make([]byte, 0, 8))
	binary.Write(buff, binary.BigEndian, uint64(i))
	binary.Read(buff, binary.BigEndian, &v)
	return v
}

func (I64) NumType() NumberType {
	return NumTypeI64
}

func (I64) ValType() ValueType {
	return ValTypeNum
}

func (I64) ValidateValueType(v types.ValueType) bool {
	if v == types.I64 {
		return true
	}
	return false
}

func (i I64) ToValue() Value {
	return i
}

func (I64) ExpectNumber() (NumberType, error) {
	return NumTypeI64, nil
}

func (i I64) ToUint64() (uint64, error) {
	buff := bytes.NewBuffer(make([]byte, 0, 8))
	if err := binary.Write(buff, binary.BigEndian, int64(i)); err != nil {
		return 0, fmt.Errorf("ToUint64: %w", err)
	}
	var v uint64
	if err := binary.Read(buff, binary.BigEndian, &v); err != nil {
		return 0, fmt.Errorf("ToUint64: %w", err)
	}
	return v, nil
}

type F32 float32

func (F32) NumType() NumberType {
	return NumTypeF32
}

func (F32) ValType() ValueType {
	return ValTypeNum
}

func (F32) ValidateValueType(v types.ValueType) bool {
	if v == types.F32 {
		return true
	}
	return false
}

func (f F32) ToValue() Value {
	return f
}

func (F32) ExpectNumber() (NumberType, error) {
	return NumTypeF32, nil
}

type F64 float64

func (F64) NumType() NumberType {
	return NumTypeF64
}

func (F64) ValType() ValueType {
	return ValTypeNum
}

func (f F64) ToValue() Value {
	return f
}

func (F64) ValidateValueType(v types.ValueType) bool {
	if v == types.F64 {
		return true
	}
	return false
}

func (F64) ExpectNumber() (NumberType, error) {
	return NumTypeF64, nil
}

type Vector [16]byte // 128bit value

func (Vector) ValType() ValueType {
	return ValTypeVec
}

func (Vector) ExpectNumber() (NumberType, error) {
	return NumberType(0xff), fmt.Errorf("Not number")
}

func Float32FromUint32(val uint32) float32 {
	return math.Float32frombits(val)
}

func Float64FromUint64(val uint64) float64 {
	return math.Float64frombits(val)
}

// unsafe
func Uint32FromFloat32(val float32) uint32 {
	return *(*uint32)(unsafe.Pointer(&val))
}

func Uint64FromFloat64(val float64) uint64 {
	return *(*uint64)(unsafe.Pointer(&val))
}

func GetNum[T NumberTypeSet](n Number) T {
	return n.(T)
}

type Result interface {
	ResultType()
}

type ResultType uint8

const (
	ResultTypeValue ResultType = 0
	ResultTypeTrap  ResultType = 1
)

type Trap struct {
	error
}

type ResultTypeSet interface {
	~int32 | ~int64 | ~float32 | ~float64 | Trap
}

func GetResult[T ResultTypeSet](r Result) T {
	return r.(T)
}

func FromString(val string, typ types.ValueType) (Value, error) {
	switch typ {
	case types.I32:
		if isNeg(val) {
			v, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return nil, err
			}
			return NewI32(int32(v)), nil
		}
		base := baseNum(val)
		v, err := strconv.ParseUint(trimBase(val, base), base, 32)
		if err != nil {
			return nil, err
		}
		return NewI32(int32(v)), nil
	case types.I64:
		if isNeg(val) {
			v, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return nil, err
			}
			return NewI64(v), nil
		}
		base := baseNum(val)
		v, err := strconv.ParseUint(trimBase(val, base), base, 64)
		if err != nil {
			return nil, err
		}
		return NewI64(v), nil
	case types.F32:
		return nil, types.NotImplemented
	case types.F64:
		return nil, types.NotImplemented
	case types.V128:
		return nil, types.NotImplemented
	default:
		return nil, types.InvalidValueType
	}
}

func isNeg(s string) bool {
	if string(s[0]) == "-" {
		return true
	}
	return false
}

func baseNum(s string) int {
	if len(s) < 2 {
		return 10
	}
	base := string(s[:2])
	switch base {
	case "0x":
		return 16
	case "0b":
		return 2
	case "0o":
		return 8
	default:
		return 10
	}
}

func trimBase(s string, base int) string {
	switch base {
	case 16, 8, 2:
		return string(s[2:])
	default:
		return s
	}
}
