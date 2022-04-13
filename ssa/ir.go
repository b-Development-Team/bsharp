package ssa

import (
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

type IRNode int

const (
	IRNodeIndex IRNode = iota
	IRNodeLength
	IRNodeMake
	IRNodeGet
	IRNodeArray
	IRNodeFn
	IRNodeExists
	IRNodeKeys
	IRNodeTime
	IRNodeGetStruct
	IRNodeCanCast
)

func (i IRNode) String() string {
	// Converts to string
	return [...]string{"Index", "Length", "Make", "Get", "Array", "Fn", "Exists", "Keys", "Time", "GetStruct", "CanCast"}[i]
}

type IRValue struct {
	Kind   IRNode
	Params []ID
	Typ    types.Type
}

func (i *IRValue) Type() types.Type { return i.Typ }
func (i *IRValue) Args() []ID       { return i.Params }
func (i *IRValue) SetArgs(v []ID)   { i.Params = v }
func (i *IRValue) String() string {
	out := &strings.Builder{}
	out.WriteString(i.Kind.String())
	out.WriteString(" (")
	for i, par := range i.Params {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(par.String())
	}
	out.WriteString(")")
	return out.String()
}

type LiveIRNode int

const (
	IRNodePrint LiveIRNode = iota
	IRNodeSet
	IRNodeAppend
	IRNodeSlice
	IRNodeSetIndex
	IRNodeSetStruct
)

func (i LiveIRNode) String() string {
	return [...]string{"Print", "Set", "Append", "Slice", "SetIndex", "SetStruct"}[i]
}

type LiveIRValue struct {
	Kind   LiveIRNode
	Params []ID
	Typ    types.Type
}

func (i *LiveIRValue) Type() types.Type { return i.Typ }
func (i *LiveIRValue) Args() []ID       { return i.Params }
func (i *LiveIRValue) SetArgs(v []ID)   { i.Params = v }
func (i *LiveIRValue) String() string {
	out := &strings.Builder{}
	out.WriteString(i.Kind.String())
	out.WriteString(" (")
	for i, par := range i.Params {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(par.String())
	}
	out.WriteString(")")
	return out.String()
}
