package optimize

import (
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (o *Optimizer) optimizePrint(c *ir.PrintNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(c.Arg)
	return &Result{
		Stmt:    ir.NewCallNode(&ir.PrintNode{Arg: v.Stmt}, pos),
		NotDead: true,
	}
}

func (o *Optimizer) optimizeConcat(c *ir.ConcatNode, pos *tokens.Pos) *Result {
	isConst := true
	stmts := make([]ir.Node, len(c.Values))
	for i, v := range c.Values {
		res := o.OptimizeNode(v)
		stmts[i] = res.Stmt
		if !res.IsConst {
			isConst = false
		}
	}

	// If is const, concat
	if isConst {
		outVal := &strings.Builder{}
		for _, v := range stmts {
			outVal.WriteString(v.(*ir.Const).Value.(string))
		}
		return &Result{
			Stmt:    ir.NewConst(types.STRING, pos, outVal.String()),
			IsConst: true,
		}
	}

	return &Result{
		Stmt: ir.NewCallNode(&ir.ConcatNode{
			Values: stmts,
		}, pos),
		IsConst: false,
	}
}
