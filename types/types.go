package types

import (
	"fmt"
	"strings"
)

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
	FUNCTION
	NULL

	// Special types
	ANY
	VARIADIC
	IDENT
)

var basicTypeNames = map[BasicType]string{
	INT:      "INT",
	FLOAT:    "FLOAT",
	STRING:   "STRING",
	BOOL:     "BOOL",
	ARRAY:    "ARRAY",
	MAP:      "MAP",
	FUNCTION: "FUNCTION",
	NULL:     "NULL",
	ANY:      "ANY",
	VARIADIC: "VARIADIC",
	IDENT:    "IDENT",
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

type MulType struct {
	Typs []Type
}

func (a *MulType) Equal(b Type) bool {
	for _, typ := range a.Typs {
		if typ.Equal(b) {
			return true
		}
	}
	return false
}

func (a *MulType) String() string {
	typs := make([]string, len(a.Typs))
	for i, typ := range a.Typs {
		typs[i] = typ.String()
	}
	return "(" + strings.Join(typs, "|") + ")"
}

func (a *MulType) BasicType() BasicType {
	return ANY
}

func NewMulType(typs ...Type) *MulType {
	return &MulType{
		Typs: typs,
	}
}

type FuncType struct {
	ParTypes []Type
	RetType  Type
}

func (f *FuncType) BasicType() BasicType {
	return FUNCTION
}

func (f *FuncType) Equal(t Type) bool {
	if t.BasicType() != FUNCTION {
		return false
	}

	v := t.(*FuncType)
	if len(v.ParTypes) != len(f.ParTypes) {
		return false
	}
	for i, par := range f.ParTypes {
		if !v.ParTypes[i].Equal(par) {
			return false
		}
	}
	return f.RetType.Equal(v.RetType)
}

func (f *FuncType) String() string {
	out := &strings.Builder{}
	out.WriteString("FUNC{")
	for i, parType := range f.ParTypes {
		out.WriteString(parType.String())
		if i != len(f.ParTypes)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(f.RetType.String())
	return out.String()
}

func NewFuncType(parTyps []Type, retType Type) Type {
	return &FuncType{
		ParTypes: parTyps,
		RetType:  retType,
	}
}
