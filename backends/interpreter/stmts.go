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

		case *ir.RandintNode:
			return i.evalRandint(c)

		case *ir.RandomNode:
			return i.evalRandom(c)

		case *ir.MathFunctionNode:
			return i.evalMathFn(n.Pos(), c)

		case *ir.MakeNode:
			return i.evalMake(n.Pos(), c)

		case *ir.SetNode:
			return i.evalSet(n.Pos(), c)

		case *ir.GetNode:
			return i.evalGet(n.Pos(), c)

		case *ir.SwitchNode:
			return NewValue(types.NULL, nil), i.evalSwitchNode(c)

		case *ir.ArrayNode:
			return i.evalArray(c)

		case *ir.AppendNode:
			return i.evalAppend(c)

		case *ir.LogicalOpNode:
			return i.evalLogicalOp(n.Pos(), c)

		case *ir.FnNode:
			return i.evalFnNode(c)

		case *ir.FnCallNode:
			return i.evalCallNode(c)

		case *ir.ExistsNode:
			return i.evalExists(n.Pos(), c)

		case *ir.KeysNode:
			return i.evalKeys(n.Pos(), c)

		default:
			return nil, n.Pos().Error("unknown call node: %T", c)
		}

	case *ir.FnCallNode:
		return i.evalCallNode(n)

	case *ir.Const:
		return i.evalConst(n)

	case *ir.ExtensionCall:
		return i.evalExtensionCall(n)

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
