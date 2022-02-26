package ir

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type Node interface {
	Type() types.Type
	Pos() *tokens.Pos
}

type Param struct {
	ID   int
	Name string
	Type types.Type
	Pos  *tokens.Pos
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
	Body     []Node
	imported map[string]empty
	currFn   string
}

func NewBuilder() *Builder {
	return &Builder{
		Scope:    NewScope(),
		Funcs:    make(map[string]*Function),
		imported: make(map[string]empty),
	}
}
