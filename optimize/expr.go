package optimize

import (
	"math"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (o *Optimizer) optimizeMath(pos *tokens.Pos, c *ir.MathNode) *Result {
	l := o.OptimizeNode(c.Lhs)
	r := o.OptimizeNode(c.Rhs)
	if l.IsConst && r.IsConst {
		switch l.Stmt.Type() {
		case types.FLOAT:
			var out float64
			lhs := l.Stmt.(*ir.Const).Value.(float64)
			rhs := r.Stmt.(*ir.Const).Value.(float64)
			switch c.Op {
			case ir.MathOperationAdd:
				out = lhs + rhs

			case ir.MathOperationSub:
				out = lhs - rhs

			case ir.MathOperationDiv:
				out = lhs / rhs

			case ir.MathOperationMul:
				out = lhs * rhs

			case ir.MathOperationPow:
				out = math.Pow(lhs, rhs)

			case ir.MathOperationMod:
				out = math.Mod(lhs, rhs)
			}
			return &Result{
				Stmt:    ir.NewConst(types.FLOAT, pos, out),
				IsConst: true,
			}

		case types.INT:
			var out int
			lhs := l.Stmt.(*ir.Const).Value.(int)
			rhs := r.Stmt.(*ir.Const).Value.(int)
			switch c.Op {
			case ir.MathOperationAdd:
				out = lhs + rhs

			case ir.MathOperationSub:
				out = lhs - rhs

			case ir.MathOperationDiv:
				out = lhs / rhs

			case ir.MathOperationMul:
				out = lhs * rhs

			case ir.MathOperationPow:
				out = lhs ^ rhs

			case ir.MathOperationMod:
				out = lhs % rhs
			}
			return &Result{
				Stmt:    ir.NewConst(types.INT, pos, out),
				IsConst: true,
			}
		}
	}

	return &Result{
		Stmt:    ir.NewCallNode(ir.NewMathNode(c.Op, l.Stmt, r.Stmt, c.Type()), pos),
		IsConst: false,
	}
}

func (o *Optimizer) optimizeCompare(c *ir.CompareNode, pos *tokens.Pos) *Result {
	l := o.OptimizeNode(c.Lhs)
	r := o.OptimizeNode(c.Rhs)
	if l.IsConst && r.IsConst {
		switch l.Stmt.Type() {
		case types.FLOAT:
			var out bool
			lhs := l.Stmt.(*ir.Const).Value.(float64)
			rhs := r.Stmt.(*ir.Const).Value.(float64)
			switch c.Op {
			case ir.CompareOperationEqual:
				out = lhs == rhs

			case ir.CompareOperationNotEqual:
				out = lhs != rhs

			case ir.CompareOperationGreater:
				out = lhs > rhs

			case ir.CompareOperationGreaterEqual:
				out = lhs >= rhs

			case ir.CompareOperationLess:
				out = lhs < rhs

			case ir.CompareOperationLessEqual:
				out = lhs <= rhs
			}
			return &Result{
				Stmt:    ir.NewConst(types.BOOL, pos, out),
				IsConst: true,
			}

		case types.INT:
			var out bool
			lhs := l.Stmt.(*ir.Const).Value.(int)
			rhs := r.Stmt.(*ir.Const).Value.(int)
			switch c.Op {
			case ir.CompareOperationEqual:
				out = lhs == rhs

			case ir.CompareOperationNotEqual:
				out = lhs != rhs

			case ir.CompareOperationGreater:
				out = lhs > rhs

			case ir.CompareOperationGreaterEqual:
				out = lhs >= rhs

			case ir.CompareOperationLess:
				out = lhs < rhs

			case ir.CompareOperationLessEqual:
				out = lhs <= rhs
			}
			return &Result{
				Stmt:    ir.NewConst(types.BOOL, pos, out),
				IsConst: true,
			}

		case types.STRING:
			var out bool
			lhs := l.Stmt.(*ir.Const).Value.(string)
			rhs := r.Stmt.(*ir.Const).Value.(string)
			switch c.Op {
			case ir.CompareOperationEqual:
				out = lhs == rhs

			case ir.CompareOperationNotEqual:
				out = lhs != rhs

			case ir.CompareOperationGreater:
				out = lhs > rhs

			case ir.CompareOperationGreaterEqual:
				out = lhs >= rhs

			case ir.CompareOperationLess:
				out = lhs < rhs

			case ir.CompareOperationLessEqual:
				out = lhs <= rhs
			}
			return &Result{
				Stmt:    ir.NewConst(types.BOOL, pos, out),
				IsConst: true,
			}
		}
	}

	return &Result{
		Stmt: ir.NewCallNode(&ir.CompareNode{
			Op:  c.Op,
			Lhs: l.Stmt,
			Rhs: r.Stmt,
		}, pos),
		IsConst: false,
	}
}
