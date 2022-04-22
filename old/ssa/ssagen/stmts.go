package ssagen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/old/ssa"
)

func (s *SSAGen) Add(node ir.Node) ssa.ID {
	if s.blk == nil {
		return ssa.NullID()
	}

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
			return s.newLiveIRNode(ssa.IRNodePrint, n)

		case *ir.IndexNode:
			return s.newIRNode(ssa.IRNodeIndex, n)

		case *ir.LengthNode:
			return s.newIRNode(ssa.IRNodeLength, n)

		case *ir.MakeNode:
			return s.newIRNode(ssa.IRNodeMake, n)

		case *ir.SetNode:
			return s.newLiveIRNode(ssa.IRNodeSet, n)

		case *ir.GetNode:
			return s.newIRNode(ssa.IRNodeGet, n)

		case *ir.ArrayNode:
			return s.newIRNode(ssa.IRNodeArray, n)

		case *ir.AppendNode:
			return s.newLiveIRNode(ssa.IRNodeAppend, n)

		case *ir.FnNode:
			return s.newIRNode(ssa.IRNodeFn, n)

		case *ir.ExistsNode:
			return s.newIRNode(ssa.IRNodeExists, n)

		case *ir.KeysNode:
			return s.newIRNode(ssa.IRNodeKeys, n)

		case *ir.TimeNode:
			return s.newIRNode(ssa.IRNodeTime, n)

		case *ir.SliceNode:
			return s.newLiveIRNode(ssa.IRNodeSlice, n)

		case *ir.SetIndexNode:
			return s.newLiveIRNode(ssa.IRNodeSetIndex, n)

		case *ir.GetStructNode:
			return s.newIRNode(ssa.IRNodeGetStruct, n)

		case *ir.SetStructNode:
			return s.newLiveIRNode(ssa.IRNodeSetStruct, n)

		case *ir.CanCastNode:
			return s.newIRNode(ssa.IRNodeCanCast, n)

		case *ir.CastNode:
			return s.blk.AddInstruction(&ssa.Cast{Value: s.Add(c.Value), From: c.Value.Type(), To: c.Type()}, n.Pos())

		case *ir.ReturnNode:
			return s.addReturn(c)

		case *ir.ConcatNode:
			return s.addConcat(n.Pos(), c)

		case *ir.LogicalOpNode:
			return s.addLogicalOp(n.Pos(), c)

		default:
			panic(fmt.Sprintf("unknown call node type: %T", c))
		}

	case *ir.BlockNode:
		switch b := n.Block.(type) {
		case *ir.IfNode:
			return s.addIf(n.Pos(), b)

		case *ir.WhileNode:
			return s.addWhile(n.Pos(), b)

		case *ir.SwitchNode:
			return s.addSwitch(n.Pos(), b)

		default:
			panic(fmt.Sprintf("unknown block node type: %T", b))
		}

	case *ir.Const:
		return s.addConst(n)

	case *ir.CastNode:
		return s.blk.AddInstruction(&ssa.Cast{Value: s.Add(n.Value), From: n.Value.Type(), To: n.Type()}, n.Pos())

	case *ir.ExtensionCall:
		return s.addExtensionCall(n)

	case *ir.FnCallNode:
		return s.addFnCall(n)

	default:
		panic(fmt.Sprintf("unknown node type: %T", n))
	}
}
