package ir

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type Node interface {
	Type() types.Type
	Pos() *tokens.Pos
}

type Call interface {
	Type() types.Type
	Args() []Node
}

type Block interface{} // TODO: Have something in this to make sure its a block

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
	pos     *tokens.Pos
	Scope   *ScopeInfo
}

func (f *Function) Pos() *tokens.Pos { return f.pos }

type empty struct{}

type ErrorLevel int

const (
	ErrorLevelError   ErrorLevel = iota
	ErrorLevelWarning            // TODO: Use this
)

type Error struct {
	Level   ErrorLevel
	Pos     *tokens.Pos
	Message string
}

type Builder struct {
	Scope      *Scope
	Funcs      map[string]*Function
	Body       []Node
	imported   map[string]empty
	currFn     string
	extensions map[string]*Extension
	typeNames  map[string]types.Type
	consts     map[string]*Const
	Errors     []*Error
}

func (b *Builder) AddExtension(e *Extension) {
	b.extensions[e.Name] = e
}

func (b *Builder) Error(level ErrorLevel, pos *tokens.Pos, format string, args ...interface{}) {
	b.Errors = append(b.Errors, &Error{Level: level, Pos: pos, Message: fmt.Sprintf(format, args...)})
}

type IR struct {
	Funcs       map[string]*Function
	Variables   []*Variable
	Body        []Node
	GlobalScope *ScopeInfo
}

func NewBuilder() *Builder {
	return &Builder{
		Scope:      NewScope(),
		Funcs:      make(map[string]*Function),
		imported:   make(map[string]empty),
		extensions: make(map[string]*Extension),
		typeNames:  make(map[string]types.Type),
		consts:     make(map[string]*Const),
	}
}

func (b *Builder) IR() *IR {
	return &IR{
		Funcs:       b.Funcs,
		Variables:   b.Scope.Variables,
		Body:        b.Body,
		GlobalScope: b.Scope.CurrScopeInfo(),
	}
}

type TypedValue struct{ typ types.Type }

func (t *TypedValue) Type() types.Type { return t.typ }

func (t *TypedValue) Args() []Node { return []Node{} }

func NewTypedValue(typ types.Type) Call { return &TypedValue{typ: typ} }

type TypedNode struct {
	typ types.Type
	pos *tokens.Pos
}

func (t *TypedNode) Type() types.Type { return t.typ }

func (t *TypedNode) Pos() *tokens.Pos { return t.pos }

func NewTypedNode(typ types.Type, pos *tokens.Pos) Node {
	return &TypedNode{typ: typ, pos: pos}
}
