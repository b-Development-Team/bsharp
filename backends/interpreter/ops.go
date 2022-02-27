package interpreter

import (
	"math"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalMathNode(pos *tokens.Pos, node *ir.MathNode) (*Value, error) {
	left, err := i.evalNode(node.Lhs)
	if err != nil {
		return nil, err
	}

	right, err := i.evalNode(node.Rhs)
	if err != nil {
		return nil, err
	}

	switch left.Type {
	case types.INT:
		v := intOp(left.Value.(int), right.Value.(int), node.Op)
		if v == nil {
			return nil, pos.Error("unknown math operation: %s", node.Op)
		}
		return NewValue(types.INT, *v), nil

	case types.FLOAT:
		v := floatOp(left.Value.(float64), right.Value.(float64), node.Op)
		if v == nil {
			return nil, pos.Error("unknown math operation: %s", node.Op)
		}
		return NewValue(types.INT, *v), nil

	default:
		return nil, pos.Error("cannot perform operations on type %s", left.Type.String())
	}
}

func intOp(lhs int, rhs int, op ir.MathOperation) *int {
	var out int
	switch op {
	case ir.MathOperationAdd:
		out = lhs + rhs

	case ir.MathOperationSub:
		out = lhs - rhs

	case ir.MathOperationMul:
		out = lhs * rhs

	case ir.MathOperationDiv:
		out = lhs / rhs

	case ir.MathOperationPow:
		out = lhs ^ rhs

	case ir.MathOperationMod:
		out = lhs % rhs

	default:
		return nil
	}

	return &out
}

func floatOp(lhs float64, rhs float64, op ir.MathOperation) *float64 {
	var out float64
	switch op {
	case ir.MathOperationAdd:
		out = lhs + rhs

	case ir.MathOperationSub:
		out = lhs - rhs

	case ir.MathOperationMul:
		out = lhs * rhs

	case ir.MathOperationDiv:
		out = lhs / rhs

	case ir.MathOperationPow:
		out = math.Pow(lhs, rhs)

	case ir.MathOperationMod:
		out = math.Mod(lhs, rhs)

	default:
		return nil
	}

	return &out
}
