package interpreter

import (
	"io"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Interpreter struct {
	ir *ir.IR

	stdout     io.Writer
	stack      *stack
	extensions map[string]*Extension

	retVal *Value
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
		extensions: make(map[string]*Extension),
		stack: &stack{
			vals: make([]scope, 0),
		},
	}
}

func (i *Interpreter) AddExtension(e *Extension) {
	i.extensions[e.Name] = e
}

func (i *Interpreter) SetStdout(stdout io.Writer) {
	i.stdout = stdout
}

func (i *Interpreter) Run() error {
	i.stack.Push()
	for _, node := range i.ir.Body {
		if _, err := i.evalNode(node); err != nil {
			return err
		}
	}
	i.stack.Pop()
	return nil
}
