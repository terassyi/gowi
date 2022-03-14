package value

type Number interface {
	NumType() NumberType
}

type Reference interface {
	RefType() ReferenceType
}

type Value interface {
	ValType() ValueType
}

type ExternalVal interface {
	ExValType() ExternValueType
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

type ExternValueType uint8

const (
	ExternValTypeFunc   ExternValueType = 0
	ExternValTypeTable  ExternValueType = 1
	ExternValTypeMem    ExternValueType = 2
	ExternValTypeGlobal ExternValueType = 3
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
