package bstar

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
)

func (b *BStar) buildFn(fn *ir.Function) (Node, error) {
	args := []Node{constNode("FUNC"), constNode(fn.Name)}
	pars := make([]Node, len(fn.Params))
	for i, arg := range fn.Params {
		pars[i] = constNode(`"` + formatName(fmt.Sprintf("%s%d", arg.Name, arg.ID)) + `"`)
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

func (b *BStar) buildFnMap() Node {
	bod := make([]Node, len(b.ir.Funcs))
	i := 0
	for _, fn := range b.ir.Funcs {
		call := make([]Node, len(fn.Params)+1)
		call[0] = constNode(fn.Name)
		for j := range fn.Params {
			call[j+1] = blockNode(true, constNode("INDEX"), blockNode(true, constNode("VAR"), constNode("args")), constNode(j))
		}
		bod[i] = blockNode(false, constNode("IF"),
			blockNode(true, constNode("COMPARE"), blockNode(true, constNode("VAR"), constNode("name")), constNode("=="), constNode(fmt.Sprintf("%q", fn.Name))),
			blockNode(false, constNode("RETURN"), blockNode(false, call...)),
			b.noPrintNode(),
		)
		i++
	}
	return blockNode(false, constNode("FUNC"), constNode("CALL"), blockNode(true, constNode("ARRAY"), constNode(`"name"`), constNode(`"args"`)), blockNode(false, append([]Node{constNode("BLOCK")}, bod...)...))
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

	fn, err := b.buildNode(n.Fn)
	if err != nil {
		return nil, err
	}
	args := make([]Node, len(n.Params)+1)
	args[0] = constNode("ARRAY")
	for i, arg := range n.Params {
		args[i+1], err = b.buildNode(arg)
		if err != nil {
			return nil, err
		}
	}

	return blockNode(true, constNode("CALL"), fn, blockNode(true, args...)), nil
}
