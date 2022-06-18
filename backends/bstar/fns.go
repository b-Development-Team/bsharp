package bstar

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BStar) buildFn(fn *ir.Function) (Node, error) {
	args := []Node{constNode("FUNC"), constNode(fn.Name)}
	pars := make([]Node, len(fn.Params))
	for i, arg := range fn.Params {
		pars[i] = constNode(fmt.Sprintf("\"%s%d\"", arg.Name, arg.ID))
	}
	args = append(args, blockNode(true, append([]Node{constNode("ARRAY")}, pars...)...))
	body, err := b.buildNodes(fn.Body)
	if err != nil {
		return nil, err
	}
	bod := []Node{constNode("BLOCK")}
	bod = append(bod, body...)
	bod = append(bod, b.noPrintNode())
	args = append(args, blockNode(false, bod...))
	return blockNode(false, args...), nil
}

func (b *BStar) buildFnCall(n *ir.FnCallNode) (Node, error) {
	v, ok := n.Fn.(*ir.CallNode)
	var c *ir.FnNode
	if ok {
		c, ok = v.Call.(*ir.FnNode)
	}
	if ok {
		args := make([]Node, len(n.Params))
		var err error
		for i, par := range n.Params {
			args[i], err = b.buildNode(par)
			if err != nil {
				return nil, err
			}
		}

		return blockNode(true, append([]Node{constNode(c.Name)}, args...)...), nil
	}

	return nil, n.Pos().Error("first-class functions not implemented yet")
}
