package ir

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type CodeConfig struct {
	Indent int
}

type Node interface {
	Type() types.Type
	Pos() *tokens.Pos
	Code(CodeConfig) string
}

type Call interface {
	Type() types.Type
	Code(CodeConfig) string
}

type Block interface {
	Code(CodeConfig) string
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
	pos     *tokens.Pos
	Scope   *ScopeInfo
}

func (f *Function) Pos() *tokens.Pos { return f.pos }
func (f *Function) Code(cnf CodeConfig) string {
	out := &strings.Builder{}
	fmt.Fprintf(out, "[FUNC %s ", f.Name)
	for _, par := range f.Params {
		fmt.Fprintf(out, "[PARAM %s %s] ", par.Name, par.Type.String())
	}
	if !types.NULL.Equal(f.RetType) {
		fmt.Fprintf(out, "[RETURNS %s]", f.RetType.String())
	}
	out.WriteString("\n")
	out.WriteString(bodyCode(cnf, f.Body))
	out.WriteString("\n]")
	return out.String()
}

type empty struct{}

type Builder struct {
	Scope      *Scope
	Funcs      map[string]*Function
	Body       []Node
	imported   map[string]empty
	currFn     string
	extensions map[string]*Extension
}

func (b *Builder) AddExtension(e *Extension) {
	b.extensions[e.Name] = e
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

func (i *IR) Code(cnf CodeConfig) string {
	out := &strings.Builder{}
	for _, fn := range i.Funcs {
		out.WriteString(fn.Code(cnf))
		out.WriteString("\n")
	}
	for _, node := range i.Body {
		out.WriteString(node.Code(cnf))
		out.WriteString("\n")
	}
	return out.String()
}
