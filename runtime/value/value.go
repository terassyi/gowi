package value

import "math"

type Number interface {
	NumType() NumberType
}

type NumberTypeSet interface {
	~int32 | ~int64 | ~float32 | ~float64
}

type Reference interface {
	RefType() ReferenceType
}

type Value interface {
	ValType() ValueType
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

type I32 int32

func (I32) NumType() NumberType {
	return NumTypeI32
}

func (I32) ValType() ValueType {
	return ValTypeNum
}

type I64 int64

func (I64) NumType() NumberType {
	return NumTypeI64
}

func (I64) ValType() ValueType {
	return ValTypeNum
}

type F32 float32

func (F32) NumType() NumberType {
	return NumTypeF32
}

func (F32) ValType() ValueType {
	return ValTypeNum
}

type F64 float64

func (F64) NumType() NumberType {
	return NumTypeF64
}

func (F64) ValType() ValueType {
	return ValTypeNum
}

type Vector [16]byte // 128bit value

func (Vector) ValType() ValueType {
	return ValTypeVec
}

func Float32FromUint32(val uint32) float32 {
	return math.Float32frombits(val)
}

func Float64FromUint64(val uint64) float64 {
	return math.Float64frombits(val)
}

func GetNum[T NumberTypeSet](n Number) T {
	return n.(T)
}
