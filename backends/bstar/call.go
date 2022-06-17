package bstar

import (
	"fmt"
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

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
		return blockNode(false, constNode("DEFINE"), constNode(v.Name+strconv.Itoa(v.ID)), args[1]), nil

	case *ir.TimeNode:
		switch c.Mode {
		case ir.TimeModeMicro:
			return blockNode(true, constNode("INT"), blockNode(true, constNode("MUL"), blockNode(true, constNode("TIME")), constNode(1000000))), nil

		case ir.TimeModeMilli:
			return blockNode(true, constNode("INT"), blockNode(true, constNode("MUL"), blockNode(true, constNode("TIME")), constNode(1000))), nil

		case ir.TimeModeNano:
			return blockNode(true, constNode("INT"), blockNode(true, constNode("MUL"), blockNode(true, constNode("TIME")), constNode(1e9))), nil
		}

		return blockNode(true, constNode("INT"), blockNode(true, constNode("TIME"))), nil

	case *ir.MathNode:
		return blockNode(true, args...), nil

	case *ir.VarNode:
		v := b.ir.Variables[c.ID]
		return blockNode(true, constNode("VAR"), constNode(v.Name+strconv.Itoa(v.ID))), nil

	case *ir.CastNode:
		switch c.Type() {
		case types.INT:
			return blockNode(true, constNode("INT"), args[0]), nil

		case types.FLOAT:
			return blockNode(true, constNode("FLOAT"), args[0]), nil

		case types.STRING:
			return blockNode(true, constNode("STR"), args[0]), nil
		}
		return args[0], nil

	case *ir.PrintNode:
		return blockNode(false, constNode("STR"), args[0]), nil

	case *ir.ConcatNode:
		return blockNode(true, append([]Node{constNode("CONCAT")}, args...)...), nil

	default:
		return nil, fmt.Errorf("unknown call: %T", c)
	}
}
