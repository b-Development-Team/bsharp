package interpreter

import "github.com/Nv7-Github/bsharp/ir"

func (i *Interpreter) evalConst(c *ir.Const) (*Value, error) {
	return NewValue(c.Type(), c.Value), nil
}
