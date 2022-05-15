package bsp

import (
	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BSP) buildNode(node ir.Node) (string, error) {
	switch n := node.(type) {
	case *ir.CallNode:
		switch c := n.Call.(type) {
		case *ir.PrintNode:
			return b.addCall("PRINT", c)

		case *ir.MathNode:
			return b.addCall("MATH", c)

		case *ir.DefineNode:
			return b.addCall("DEFINE", c)

		case *ir.TimeNode:
			return b.addCall("TIME", c)

		case *ir.FnNode:
			return b.addCall("FN", c)

		case *ir.VarNode:
			return b.addCall("VAR", c)

		case *ir.ConcatNode:
			return b.addCall("CONCAT", c)

		case *ir.CompareNode:
			return b.addCall("COMPARE", c)

		case *ir.IndexNode:
			return b.addCall("INDEX", c)

		case *ir.LengthNode:
			return b.addCall("LENGTH", c)

		case *ir.MakeNode:
			return b.addCall("MAKE", c)

		case *ir.SetNode:
			return b.addCall("SET", c)

		case *ir.GetNode:
			return b.addCall("GET", c)

		case *ir.ArrayNode:
			return b.addCall("ARRAY", c)

		case *ir.AppendNode:
			return b.addCall("APPEND", c)

		case *ir.ReturnNode:
			return b.addCall("RETURN", c)

		case *ir.CastNode:
			return b.addCast(c)

		case *ir.FnCallNode:
			return b.addFnCall(c)

		default:
			return "", node.Pos().Error("unsupported call type: %T", c)
		}

	case *ir.Const:
		return b.addConst(n)

	case *ir.FnCallNode:
		return b.addFnCall(n)

	case *ir.CastNode:
		return b.addCast(n)

	case *ir.BlockNode:
		switch bl := n.Block.(type) {
		case *ir.IfNode:
			return b.buildIf(bl)

		case *ir.WhileNode:
			return b.buildWhile(bl)

		case *ir.SwitchNode:
			return b.buildSwitch(bl)

		default:
			return "", node.Pos().Error("unsupported block type: %T", bl)
		}

	default:
		return "", n.Pos().Error("unknown node type: %T", node)
	}
}
