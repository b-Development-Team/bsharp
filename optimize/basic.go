package optimize

import (
	"math"
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

func (o *Optimizer) optimizeRandint(c *ir.RandintNode, pos *tokens.Pos) *Result {
	low := o.OptimizeNode(c.Lower)
	up := o.OptimizeNode(c.Upper)
	return &Result{
		Stmt: ir.NewCallNode(&ir.RandintNode{
			Lower: low.Stmt,
			Upper: up.Stmt,
		}, pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeRandom(c *ir.RandomNode, pos *tokens.Pos) *Result {
	low := o.OptimizeNode(c.Lower)
	up := o.OptimizeNode(c.Upper)
	return &Result{
		Stmt: ir.NewCallNode(&ir.RandomNode{
			Lower: low.Stmt,
			Upper: up.Stmt,
		}, pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeMathFunction(n *ir.MathFunctionNode, pos *tokens.Pos) *Result {
	v := o.OptimizeNode(n.Arg)
	if v.IsConst {
		in := v.Stmt.(*ir.Const).Value
		var out interface{}
		switch n.Func {
		case ir.MathFunctionCeil:
			out = int(math.Ceil(in.(float64)))

		case ir.MathFunctionFloor:
			out = int(math.Floor(in.(float64)))

		case ir.MathFunctionRound:
			out = int(math.Round(in.(float64)))
		}
		return &Result{
			Stmt:    ir.NewConst(n.Type(), pos, out),
			IsConst: true,
		}
	}
	return &Result{
		Stmt:    ir.NewCallNode(ir.NewMathFunctionNode(n.Func, v.Stmt, n.Type()), pos),
		IsConst: false,
	}
}
