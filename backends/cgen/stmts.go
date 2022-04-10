package cgen

import (
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
)

func (cg *CGen) AddNode(node ir.Node) (*Code, error) {
	if cg.isReturn {
		cg.isReturn = false
	}
	switch n := node.(type) {
	case *ir.CallNode:
		switch c := n.Call.(type) {
		case *ir.PrintNode:
			return cg.addPrint(c)

		case *ir.FnCallNode:
			return cg.addFnCall(c)

		case *ir.FnNode:
			return cg.addFn(c), nil

		case *ir.DefineNode:
			return cg.addDefine(c)

		case *ir.VarNode:
			return &Code{
				Value: Namespace + cg.ir.Variables[c.ID].Name + strconv.Itoa(c.ID),
			}, nil

		case *ir.ReturnNode:
			return cg.addReturn(c)

		case *ir.MathNode:
			return cg.addMath(c)

		case *ir.CastNode:
			return cg.addCast(c)

		case *ir.CompareNode:
			return cg.addCompare(c)

		case *ir.TimeNode:
			return cg.addTime(c), nil

		case *ir.ConcatNode:
			return cg.addConcat(c)

		case *ir.ArrayNode:
			return cg.addArray(c)

		case *ir.IndexNode:
			return cg.addIndex(n.Pos(), c)

		case *ir.AppendNode:
			return cg.addAppend(c)

		case *ir.MakeNode:
			return cg.addMake(c)

		case *ir.SetNode:
			return cg.addSet(c)

		case *ir.GetNode:
			return cg.addGet(c)

		case *ir.LengthNode:
			return cg.addLength(c)

		case *ir.KeysNode:
			return cg.addKeys(c)

		case *ir.ExistsNode:
			return cg.addExists(c)

		case *ir.LogicalOpNode:
			return cg.addLogicalOp(c)

		case *ir.SliceNode:
			return cg.addSlice(c)

		case *ir.SetIndexNode:
			return cg.addSetIndex(n.Pos(), c)

		case *ir.GetStructNode:
			return cg.addGetStruct(c)

		case *ir.SetStructNode:
			return cg.addSetStruct(c)

		case *ir.CanCastNode:
			return cg.addCanCast(c)

		default:
			return nil, n.Pos().Error("unknown call node: %T", c)
		}

	case *ir.BlockNode:
		switch b := n.Block.(type) {
		case *ir.IfNode:
			return cg.addIf(b)

		case *ir.WhileNode:
			return cg.addWhile(b)

		case *ir.SwitchNode:
			return cg.addSwitch(b)

		default:
			return nil, n.Pos().Error("unknown block node: %T", b)
		}

	case *ir.Const:
		return cg.addConst(n), nil

	case *ir.FnCallNode:
		return cg.addFnCall(n)

	case *ir.CastNode:
		return cg.addCast(n)

	default:
		return nil, n.Pos().Error("unknown node: %T", node)
	}
}
