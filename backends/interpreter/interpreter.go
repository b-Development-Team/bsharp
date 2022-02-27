package interpreter

import (
	"io"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Interpreter struct {
	ir *ir.IR

	stdout    io.Writer
	Variables []*Value

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
		ir:        ir,
		stdout:    stdout,
		Variables: make([]*Value, len(ir.Variables)),
	}
}

func (i *Interpreter) SetStdout(stdout io.Writer) {
	i.stdout = stdout
}

func (i *Interpreter) Run() error {
	for _, node := range i.ir.Body {
		if _, err := i.evalNode(node); err != nil {
			return err
		}
	}
	return nil
}
