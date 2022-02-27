package interpreter

import (
	"io"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Interpreter struct {
	ir *ir.IR

	stdout     io.Writer
	Variables  []*Value
	extensions map[string]*Extension

	retVal *Value
	scope  *Scope
}

type Value struct {
	Type  types.Type
	Value interface{}
}

func NewValue(typ types.Type, val interface{}) *Value {
	return &Value{
		Type:  typ,
		Value: val,
	}
}

func NewInterpreter(ir *ir.IR, stdout io.Writer) *Interpreter {
	return &Interpreter{
		ir:         ir,
		stdout:     stdout,
		Variables:  make([]*Value, len(ir.Variables)),
		scope:      NewScope(),
		extensions: make(map[string]*Extension),
	}
}

func (i *Interpreter) AddExtension(e *Extension) {
	i.extensions[e.Name] = e
}

func (i *Interpreter) SetStdout(stdout io.Writer) {
	i.stdout = stdout
}

func (i *Interpreter) Run() error {
	i.scope.Push()
	for _, node := range i.ir.Body {
		if _, err := i.evalNode(node); err != nil {
			return err
		}
	}
	i.pop()
	return nil
}

func (i *Interpreter) pop() {
	vals := i.scope.Pop()
	for _, v := range vals {
		i.Variables[v] = nil
	}
}
