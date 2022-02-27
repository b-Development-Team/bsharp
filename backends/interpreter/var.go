package interpreter

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalVarNode(n *ir.VarNode) (*Value, error) {
	return i.Variables[n.ID], nil
}

func (i *Interpreter) evalDefineNode(n *ir.DefineNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	i.Variables[n.Var] = v
	return NewValue(types.NULL, nil), nil
}
