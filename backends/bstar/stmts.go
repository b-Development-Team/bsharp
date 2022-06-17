package bstar

import (
	"fmt"
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BStar) buildNode(node ir.Node) (Node, error) {
	switch n := node.(type) {
	case *ir.CallNode:
		return b.buildCall(n)

	default:
		return nil, fmt.Errorf("unknown node: %T", n)
	}
}

func (b *BStar) buildCall(n *ir.CallNode) (Node, error) {
	a := n.Call.Args()
	args := make([]Node, len(a))
	for i, v := range a {
		arg, err := b.buildNode(v)
		if err != nil {
			return nil, err
		}
		args[i] = arg
	}
	switch c := n.Call.(type) {
	case *ir.DefineNode:
		v := b.ir.Variables[c.Var]
		return blockNode(constNode("DEFINE"), constNode(v.Name+strconv.Itoa(v.ID)), args[1]), nil

	default:
		return nil, fmt.Errorf("unknown call: %T", c)
	}
}
