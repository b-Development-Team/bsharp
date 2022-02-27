package interpreter

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalNode(node ir.Node) (*Value, error) {
	// If returned, ignore
	if i.retVal != nil {
		return nil, nil
	}

	switch n := node.(type) {
	case *ir.CallNode:
		switch c := n.Call.(type) {
		case *ir.PrintNode:
			return NewValue(types.NULL, nil), i.evalPrint(c)

		case *ir.CastNode:
			return i.evalCast(c)

		case *ir.ReturnNode:
			return i.evalReturnNode(c)

		case *ir.MathNode:
			return i.evalMathNode(n.Pos(), c)

		case *ir.VarNode:
			return i.evalVarNode(c)

		case *ir.DefineNode:
			return i.evalDefineNode(c)

		case *ir.ConcatNode:
			return i.evalConcat(c)

		case *ir.IfNode:
			return NewValue(types.NULL, nil), i.evalIfNode(c)

		case *ir.CompareNode:
			return i.evalCompareNode(n.Pos(), c)

		case *ir.IndexNode:
			return i.evalIndex(n.Pos(), c)

		case *ir.LengthNode:
			return i.evalLength(n.Pos(), c)

		case *ir.WhileNode:
			return NewValue(types.NULL, nil), i.evalWhileNode(c)

		default:
			return nil, n.Pos().Error("unknown call node: %T", c)
		}

	case *ir.FnCallNode:
		return i.evalCallNode(n)

	case *ir.Const:
		return i.evalConst(n)

	default:
		return nil, n.Pos().Error("unknown node type: %T", node)
	}
}

func (i *Interpreter) evalNodes(nodes []ir.Node) ([]*Value, error) {
	out := make([]*Value, len(nodes))
	for ind, node := range nodes {
		v, err := i.evalNode(node)
		if err != nil {
			return nil, err
		}
		out[ind] = v
	}
	return out, nil
}
