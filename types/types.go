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

func ParseType(typ string, names map[string]Type) (Type, error) {
	tokens, err := tokenize([]rune(typ))
	if err != nil {
		return nil, err
	}
	out, _, err := parse(tokens, names)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type BasicType int

const (
	INT BasicType = iota
	BYTE
	FLOAT
	BOOL
	STRING
	ARRAY
	MAP
	FUNCTION
	STRUCT
	NULL
	ANY

	// Special types
	ALL
	VARIADIC
	IDENT
	INVALID
)

var basicTypeNames = map[BasicType]string{
	INT:      "INT",
	BYTE:     "BYTE",
	FLOAT:    "FLOAT",
	STRING:   "STRING",
	BOOL:     "BOOL",
	ARRAY:    "ARRAY",
	MAP:      "MAP",
	FUNCTION: "FUNCTION",
	NULL:     "NIL",
	ALL:      "ALL",
	ANY:      "ANY",
	VARIADIC: "VARIADIC",
	IDENT:    "IDENT",
	STRUCT:   "STRUCT",
	INVALID:  "INVALID",
}

func (b BasicType) BasicType() BasicType {
	return b
}

func (b BasicType) String() string {
	return basicTypeNames[b]
}

func (b BasicType) Equal(t Type) bool {
	if b == ALL || t == ALL {
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
	if b.BasicType() == ALL {
		return true
	}
	if b.BasicType() != ARRAY {
		return false
	}

	return a.ElemType.Equal(b.(*ArrayType).ElemType)
}

func (a *ArrayType) String() string {
	return fmt.Sprintf("ARRAY{%s}", a.ElemType.String())
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
	if b.BasicType() == ALL {
		return true
	}
	if b.BasicType() != MAP {
		return false
	}

	return m.KeyType.Equal(b.(*MapType).KeyType) && m.ValType.Equal(b.(*MapType).ValType)
}

func (m *MapType) String() string {
	return fmt.Sprintf("MAP{%s, %s}", m.KeyType.String(), m.ValType.String())
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
	return ALL
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
	if t.BasicType() == ALL {
		return true
	}
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
	out.WriteString("}")
	out.WriteString(f.RetType.String())
	return out.String()
}

func NewFuncType(parTyps []Type, retType Type) Type {
	return &FuncType{
		ParTypes: parTyps,
		RetType:  retType,
	}
}

type StructField struct {
	Name string
	Type Type
}

type StructType struct {
	Fields []StructField
}

func (s *StructType) BasicType() BasicType {
	return STRUCT
}

func (s *StructType) Equal(t Type) bool {
	if t.BasicType() != STRUCT {
		return false
	}

	v := t.(*StructType)
	if len(v.Fields) != len(s.Fields) {
		return false
	}
	for i, field := range s.Fields {
		if !v.Fields[i].Type.Equal(field.Type) || field.Name != v.Fields[i].Name {
			return false
		}
	}
	return true
}

func (s *StructType) String() string {
	out := &strings.Builder{}
	out.WriteString("STRUCT{")
	for i, field := range s.Fields {
		fmt.Fprintf(out, "%s:%s", field.Name, field.Type.String())
		if i != len(s.Fields)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString("}")
	return out.String()
}

func NewStructField(name string, typ Type) StructField {
	return StructField{
		Name: name,
		Type: typ,
	}
}

func NewStruct(fields ...StructField) *StructType {
	return &StructType{
		Fields: fields,
	}
}
