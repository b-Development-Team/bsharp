package optimize

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (o *Optimizer) optimizeDefine(n *ir.DefineNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(n.Value)
	o.scope.SetVar(n.Var, v)
	return &Result{
		Stmt:    ir.NewCallNode(ir.NewDefineNode(n.Var, v.Stmt, o.ir), pos),
		NotDead: true,
	}
}

func (o *Optimizer) optimizeVar(n *ir.VarNode, pos *tokens.Pos) *Result {
	v := o.scope.GetVar(n.ID)
	if v != nil {
		return v
	}
	return &Result{
		Stmt:    ir.NewCallNode(n, pos),
		IsConst: false,
	}
}
