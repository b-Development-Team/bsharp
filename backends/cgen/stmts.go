package cgen

import "github.com/Nv7-Github/bsharp/ir"

func (cg *CGen) AddNode(node ir.Node) (*Code, error) {
	switch n := node.(type) {
	case *ir.CallNode:
		switch c := n.Call.(type) {
		case *ir.PrintNode:
			return cg.addPrint(c)

		case *ir.FnCallNode:
			return cg.addFnCall(c)

		case *ir.FnNode:
			return cg.addFn(c), nil

		default:
			return nil, n.Pos().Error("unknown call node: %T", c)
		}

	case *ir.Const:
		return cg.addConst(n), nil

	case *ir.FnCallNode:
		return cg.addFnCall(n)

	default:
		return nil, n.Pos().Error("unknown node: %T", node)
	}
}
