package ir

import (
	"go/token"

	"github.com/Nv7-Github/bsharp/types"
)

type Node interface {
	Type() types.Type
	Pos() *token.Pos
}

type Param struct {
	ID   int
	Name string
	Type types.Type
}

type Function struct {
	Name    string
	Params  []*Param
	RetType types.Type
	Body    []Node
}

type empty struct{}

type Builder struct {
	Scope    *Scope
	Funcs    map[string]*Function
	imported map[string]empty
}

func NewBuilder() *Builder {
	return &Builder{
		Scope:    NewScope(),
		Funcs:    make(map[string]*Function),
		imported: make(map[string]empty),
	}
}
