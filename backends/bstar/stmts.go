package bstar

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BStar) buildNode(node ir.Node) (Node, error) {
	switch n := node.(type) {
	case *ir.CallNode:
		return b.buildCall(n)

	case *ir.Const:
		return constNode(n.Value), nil

	case *ir.BlockNode:
		return b.buildBlock(n)

	default:
		return nil, fmt.Errorf("unknown node: %T", n)
	}
}
