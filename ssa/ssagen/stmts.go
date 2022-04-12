package ssagen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/ssa"
)

func (s *SSAGen) Add(node ir.Node) ssa.ID {
	switch n := node.(type) {
	case *ir.CallNode:
		switch c := n.Call.(type) {
		case *ir.DefineNode:
			return s.addDefine(n.Pos(), c)

		case *ir.VarNode:
			return s.addVar(n.Pos(), c)

		case *ir.CompareNode:
			return s.addCompare(n.Pos(), c)

		case *ir.MathNode:
			return s.addMath(n.Pos(), c)

		case *ir.PrintNode:
			return s.blk.AddInstruction(&ssa.Print{Value: s.Add(c.Arg)}, n.Pos())

		case *ir.CastNode:
			return s.blk.AddInstruction(&ssa.Cast{Value: s.Add(c.Value), From: c.Value.Type(), To: c.Type()}, n.Pos())

		default:
			panic(fmt.Sprintf("unknown call node type: %T", c))
		}

	case *ir.BlockNode:
		switch b := n.Block.(type) {
		case *ir.IfNode:
			return s.addIf(b)

		case *ir.WhileNode:
			return s.addWhile(b)

		default:
			panic(fmt.Sprintf("unknown block node type: %T", b))
		}

	case *ir.Const:
		return s.addConst(n)

	case *ir.CastNode:
		return s.blk.AddInstruction(&ssa.Cast{Value: s.Add(n.Value), From: n.Value.Type(), To: n.Type()}, n.Pos())

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}
