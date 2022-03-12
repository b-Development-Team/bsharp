package optimize

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
)

func (o *Optimizer) OptimizeNode(node ir.Node) *Result {
	switch n := node.(type) {
	case *ir.CallNode:
		switch c := n.Call.(type) {
		case *ir.MathNode:
			return o.optimizeMath(n.Pos(), c)

		case *ir.PrintNode:
			return o.optimizePrint(c, n.Pos())

		case *ir.CastNode:
			return o.optimizeCast(c, n.Pos())

		case *ir.DefineNode:
			return o.optimizeDefine(c, n.Pos())

		case *ir.VarNode:
			return o.optimizeVar(c, n.Pos())

		case *ir.WhileNode:
			return o.optimizeWhile(c, n.Pos())

		case *ir.CompareNode:
			return o.optimizeCompare(c, n.Pos())

		case *ir.ConcatNode:
			return o.optimizeConcat(c, n.Pos())

		case *ir.IfNode:
			return o.optimizeIf(c, n.Pos())

		case *ir.IndexNode:
			return o.optimizeIndex(c, n.Pos())

		case *ir.LengthNode:
			return o.optimizeLength(c, n.Pos())

		case *ir.RandintNode:
			return o.optimizeRandint(c, n.Pos())

		case *ir.RandomNode:
			return o.optimizeRandom(c, n.Pos())

		case *ir.MathFunctionNode:
			return o.optimizeMathFunction(c, n.Pos())

		case *ir.MakeNode:
			return o.optimizeMake(c, n.Pos())

		case *ir.SetNode:
			return o.optimizeSet(c, n.Pos())

		case *ir.GetNode:
			return o.optimizeGet(c, n.Pos())

		case *ir.ArrayNode:
			return o.optimizeArray(c, n.Pos())

		case *ir.AppendNode:
			return o.optimizeAppend(c, n.Pos())

		case *ir.SwitchNode:
			return o.optimizeSwitch(c, n.Pos())

		default:
			panic(fmt.Errorf("optimize: unknown call node type: %T", c)) // This shouldn't happen
		}

	case *ir.Const:
		return &Result{
			Stmt:    n,
			IsConst: true,
		}

	default:
		panic(fmt.Errorf("optimize: unknown node type: %T", n)) // This shouldn't happen
	}
}
