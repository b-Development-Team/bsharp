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
