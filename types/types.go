package types

import "fmt"

type Type interface {
	fmt.Stringer

	BasicType() BasicType
	Equal(Type) bool
}

type BasicType int

const (
	INT BasicType = iota
	FLOAT
	BOOL
	STRING
	ARRAY
	MAP
	NULL

	// Special types
	ANY
	VARIADIC
)

var basicTypeNames = map[BasicType]string{
	INT:    "int",
	FLOAT:  "float",
	STRING: "string",
	BOOL:   "bool",
	ARRAY:  "array",
	MAP:    "map",
	NULL:   "null",
}

func (b BasicType) BasicType() BasicType {
	return b
}

func (b BasicType) String() string {
	return basicTypeNames[b]
}

func (b BasicType) Equal(t Type) bool {
	if b == ANY || t == ANY {
		return true
	}
	return b == t.BasicType()
}

type ArrayType struct {
	ElemType Type
}

func NewArrayType(elemType Type) *ArrayType {
	return &ArrayType{
		ElemType: elemType,
	}
}

func (a *ArrayType) BasicType() BasicType {
	return ARRAY
}

func (a *ArrayType) Equal(b Type) bool {
	if b.BasicType() == ANY {
		return true
	}
	if b.BasicType() != ARRAY {
		return false
	}

	return a.ElemType.Equal(b.(*ArrayType).ElemType)
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("[]%s", a.ElemType.String())
}

type MapType struct {
	KeyType Type
	ValType Type
}

func NewMapType(keyType, valType Type) *MapType {
	return &MapType{
		KeyType: keyType,
		ValType: valType,
	}
}

func (m *MapType) BasicType() BasicType {
	return MAP
}

func (m *MapType) Equal(b Type) bool {
	if b.BasicType() == ANY {
		return true
	}
	if b.BasicType() != MAP {
		return false
	}

	return m.KeyType.Equal(b.(*MapType).KeyType) && m.ValType.Equal(b.(*MapType).ValType)
}

func (m *MapType) String() string {
	return fmt.Sprintf("{%s, %s}", m.KeyType.String(), m.ValType.String())
}
