package bstar

import (
	"fmt"
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
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
		return blockNode(true, append([]Node{constNode("MATH")}, args...)...), nil

	case *ir.VarNode:
		v := b.ir.Variables[c.ID]
		return blockNode(true, constNode("VAR"), constNode(v.Name+strconv.Itoa(v.ID))), nil

	case *ir.CastNode:
		return b.addCast(c)

	case *ir.PrintNode:
		return blockNode(false, constNode("CONCAT"), args[0], constNode(`"\n"`)), nil

	case *ir.ConcatNode:
		return blockNode(true, append([]Node{constNode("CONCAT")}, args...)...), nil

	case *ir.IndexNode:
		return blockNode(true, constNode("INDEX"), args[0], args[1]), nil

	case *ir.LengthNode:
		return blockNode(true, constNode("LENGTH"), args[0]), nil

	case *ir.CompareNode:
		op := c.Op.String()
		if c.Op == ir.CompareOperationEqual {
			op = "="
		}
		return blockNode(true, constNode("COMPARE"), args[0], constNode(op), args[2]), nil

	case *ir.ReturnNode:
		return blockNode(false, constNode("RETURN"), args[0]), nil

	case *ir.ArrayNode:
		return blockNode(false, append([]Node{constNode("ARRAY")}, args...)...), nil

	case *ir.AppendNode:
		cv, ok := c.Array.(*ir.CallNode)
		if !ok {
			return nil, n.Pos().Error("APPEND must be called on a VAR node in the B* backend")
		}
		v, ok := cv.Call.(*ir.VarNode)
		if !ok {
			return nil, n.Pos().Error("APPEND must be called on a VAR node in the B* backend")
		}
		va := b.ir.Variables[v.ID]
		return blockNode(false, constNode("DEFINE"), constNode(fmt.Sprintf("%s%d", va.Name, va.ID)), blockNode(true, constNode("CONCAT"), args[0], blockNode(true, constNode("ARRAY"), args[1]))), nil

	case *ir.MakeNode:
		return b.makeVal(n.Pos(), c.Type())

	case *ir.SetIndexNode:
		cv, ok := c.Array.(*ir.CallNode)
		if !ok {
			return nil, n.Pos().Error("SET must be called on a VAR node in the B* backend")
		}
		v, ok := cv.Call.(*ir.VarNode)
		if !ok {
			return nil, n.Pos().Error("SET must be called on a VAR node in the B* backend")
		}
		va := b.ir.Variables[v.ID]
		return blockNode(false, constNode("DEFINE"), constNode(fmt.Sprintf("%s%d", va.Name, va.ID)), blockNode(true, constNode("SETINDEX"), args[0], args[1], args[2])), nil

	case *ir.SetNode:
		cv, ok := c.Map.(*ir.CallNode)
		if !ok {
			return nil, n.Pos().Error("SET must be called on a VAR node in the B* backend")
		}
		v, ok := cv.Call.(*ir.VarNode)
		if !ok {
			return nil, n.Pos().Error("SET must be called on a VAR node in the B* backend")
		}
		va := b.ir.Variables[v.ID]
		return blockNode(false, constNode("DEFINE"), constNode(fmt.Sprintf("%s%d", va.Name, va.ID)), blockNode(true, constNode("MAP_SET"), args[0], args[1], args[2])), nil

	case *ir.SliceNode:
		if types.STRING.Equal(c.Value.Type()) {
			return blockNode(true, constNode("SLICE"), args[0], args[1], args[2]), nil
		}

		cv, ok := c.Value.(*ir.CallNode)
		if !ok {
			return nil, n.Pos().Error("SLICE must be called on a VAR node in the B* backend")
		}
		v, ok := cv.Call.(*ir.VarNode)
		if !ok {
			return nil, n.Pos().Error("SLICE must be called on a VAR node in the B* backend")
		}
		va := b.ir.Variables[v.ID]
		return blockNode(false, constNode("DEFINE"), constNode(fmt.Sprintf("%s%d", va.Name, va.ID)), blockNode(true, constNode("SLICE"), args[0], args[1], args[2])), nil

	case *ir.SetStructNode:
		cv, ok := c.Struct.(*ir.CallNode)
		if !ok {
			return nil, n.Pos().Error("SET must be called on a VAR node in the B* backend")
		}
		v, ok := cv.Call.(*ir.VarNode)
		if !ok {
			return nil, n.Pos().Error("SET must be called on a VAR node in the B* backend")
		}
		va := b.ir.Variables[v.ID]
		return blockNode(false, constNode("DEFINE"), constNode(fmt.Sprintf("%s%d", va.Name, va.ID)), blockNode(true, constNode("SETINDEX"), args[0], constNode(c.Field), args[2])), nil

	case *ir.GetStructNode:
		return blockNode(true, constNode("INDEX"), args[0], constNode(c.Field)), nil

	case *ir.GetNode:
		return blockNode(true, constNode("MAP_GET"), args[0], args[1]), nil

	case *ir.FnCallNode:
		return b.buildFnCall(c)

	case *ir.FnNode:
		return constNode(fmt.Sprintf("%q", c.Name)), nil

	case *ir.ExistsNode:
		// [COMPARE [FIND [INDEX [VAR x] 0] [VAR val]] != -1]
		return blockNode(true, constNode("COMPARE"),
			blockNode(true, constNode("FIND"), blockNode(true, constNode("INDEX"), args[0], constNode(0)), args[1]),
			constNode("!="),
			constNode(-1),
		), nil

	case *ir.CanCastNode:
		return blockNode(true, constNode("COMPARE"), blockNode(true, constNode("INDEX"), args[0], constNode(0)), constNode("=="), constNode(fmt.Sprintf("%q", c.Typ.String()))), nil

	case *ir.PanicNode:
		return blockNode(false, constNode("RAISE"), args[0]), nil

	case *ir.LogicalOpNode:
		switch c.Op { // 0 is False, 1 is True
		case ir.LogicalOpAnd:
			return blockNode(true, constNode("COMPARE"), blockNode(true, constNode("MATH"), args[0], constNode("+"), args[1]), constNode("=="), constNode(2)), nil // x + y == 2

		case ir.LogicalOpOr:
			return blockNode(true, constNode("COMPARE"), blockNode(true, constNode("MATH"), args[0], constNode("+"), args[1]), constNode(">"), constNode(0)), nil // x + y > 0

		case ir.LogicalOpNot:
			return blockNode(true, constNode("MATH"), constNode(1), constNode("-"), args[0]), nil // 1 - x

		default:
			return nil, n.Pos().Error("unsupported logical operator")
		}

	default:
		return nil, n.Pos().Error("unknown call: %T", c)
	}
}

func (b *BStar) addCast(c *ir.CastNode) (Node, error) {
	v, err := b.buildNode(c.Value)
	if err != nil {
		return nil, err
	}

	if types.ANY.Equal(c.Value.Type()) {
		return blockNode(true, constNode("INDEX"), v, constNode(1)), nil
	}

	switch c.Type() {
	case types.INT:
		return blockNode(true, constNode("INT"), v), nil

	case types.FLOAT:
		return blockNode(true, constNode("FLOAT"), v), nil

	case types.STRING:
		return blockNode(true, constNode("STR"), v), nil

	case types.ANY:
		return blockNode(true, constNode("ARRAY"), constNode(fmt.Sprintf("%q", c.Value.Type().String())), v), nil
	}

	return v, nil
}

func (b *BStar) makeVal(pos *tokens.Pos, t types.Type) (Node, error) {
	switch t.BasicType() {
	case types.INT:
		return constNode(0), nil

	case types.FLOAT:
		return constNode(0.0), nil

	case types.STRING:
		return constNode(`""`), nil

	case types.BOOL:
		return constNode(0), nil

	case types.ARRAY:
		return constNode("{}"), nil

	case types.MAP:
		return blockNode(true, constNode("ARRAY"), constNode("{}"), constNode("{}")), nil

	case types.STRUCT:
		t := t.(*types.StructType)
		vals := make([]Node, len(t.Fields)+1)
		vals[0] = constNode("ARRAY")
		var err error
		for i, f := range t.Fields {
			vals[i+1], err = b.makeVal(pos, f.Type)
			if err != nil {
				return nil, err
			}
		}
		return blockNode(true, vals...), err

	case types.ANY:
		return blockNode(true, constNode("ARRAY"), constNode("ANY"), constNode(-1)), nil

	default:
		return nil, pos.Error("invalid MAKE type")
	}
}
