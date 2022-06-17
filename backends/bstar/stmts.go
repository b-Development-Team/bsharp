package bstar

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (b *BStar) buildNode(node ir.Node) (Node, error) {
	switch n := node.(type) {
	case *ir.CallNode:
		return b.buildCall(n)

	case *ir.Const:
		if n.Type().Equal(types.STRING) {
			return constNode(fmt.Sprintf("%q", n.Value)), nil
		}
		return constNode(n.Value), nil

	case *ir.BlockNode:
		return b.buildBlock(n)

	default:
		return nil, n.Pos().Error("unknown node: %T", n)
	}
}
