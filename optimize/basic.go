package optimize

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
)

func (o *Optimizer) optimizePrint(c *ir.PrintNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(c.Arg)
	return &Result{
		Stmt:    ir.NewCallNode(&ir.PrintNode{Arg: v.Stmt}, pos),
		NotDead: true,
	}
}
