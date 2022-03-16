package interpreter

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalVarNode(n *ir.VarNode) (*Value, error) {
	return i.stack.Get(n.ID), nil
}

func (i *Interpreter) evalDefineNode(n *ir.DefineNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	i.stack.Set(n.Var, v, n.InScope)
	return NewValue(types.NULL, nil), nil
}
