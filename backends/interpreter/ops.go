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
		v := intMathOp(left.Value.(int), right.Value.(int), node.Op)
		if v == nil {
			return nil, pos.Error("unknown math operation: %s", node.Op)
		}
		return NewValue(types.INT, *v), nil

	case types.FLOAT:
		v := floatMathOp(left.Value.(float64), right.Value.(float64), node.Op)
		if v == nil {
			return nil, pos.Error("unknown math operation: %s", node.Op)
		}
		return NewValue(types.INT, *v), nil

	default:
		return nil, pos.Error("cannot perform operations on type %s", left.Type.String())
	}
}

func intMathOp(lhs int, rhs int, op ir.MathOperation) *int {
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

func floatMathOp(lhs float64, rhs float64, op ir.MathOperation) *float64 {
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

func (i *Interpreter) evalCompareNode(pos *tokens.Pos, node *ir.CompareNode) (*Value, error) {
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
		v := intCompOp(left.Value.(int), right.Value.(int), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	case types.FLOAT:
		v := floatCompOp(left.Value.(float64), right.Value.(float64), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	case types.STRING:
		v := stringCompOp(left.Value.(string), right.Value.(string), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	default:
		return nil, pos.Error("cannot perform comparison on type %s", left.Type.String())
	}
}

func intCompOp(lhs int, rhs int, op ir.CompareOperation) *bool {
	var out bool
	switch op {
	case ir.CompareOperationEqual:
		out = lhs == rhs

	case ir.CompareOperationNotEqual:
		out = lhs != rhs

	case ir.CompareOperationGreater:
		out = lhs > rhs

	case ir.CompareOperationGreaterEqual:
		out = lhs >= rhs

	case ir.CompareOperationLess:
		out = lhs < rhs

	case ir.CompareOperationLessEqual:
		out = lhs <= rhs

	default:
		return nil
	}

	return &out
}

func floatCompOp(lhs float64, rhs float64, op ir.CompareOperation) *bool {
	var out bool
	switch op {
	case ir.CompareOperationEqual:
		out = lhs == rhs

	case ir.CompareOperationNotEqual:
		out = lhs != rhs

	case ir.CompareOperationGreater:
		out = lhs > rhs

	case ir.CompareOperationGreaterEqual:
		out = lhs >= rhs

	case ir.CompareOperationLess:
		out = lhs < rhs

	case ir.CompareOperationLessEqual:
		out = lhs <= rhs

	default:
		return nil
	}

	return &out
}

func stringCompOp(lhs string, rhs string, op ir.CompareOperation) *bool {
	var out bool
	switch op {
	case ir.CompareOperationEqual:
		out = lhs == rhs

	case ir.CompareOperationNotEqual:
		out = lhs != rhs

	case ir.CompareOperationGreater:
		out = lhs > rhs

	case ir.CompareOperationGreaterEqual:
		out = lhs >= rhs

	case ir.CompareOperationLess:
		out = lhs < rhs

	case ir.CompareOperationLessEqual:
		out = lhs <= rhs

	default:
		return nil
	}

	return &out
}

func (i *Interpreter) evalLogicalOp(pos *tokens.Pos, node *ir.LogicalOpNode) (*Value, error) {
	val, err := i.evalNode(node.Val)
	if err != nil {
		return nil, err
	}

	var right *Value
	if node.Rhs != nil {
		right, err = i.evalNode(node.Rhs)
		if err != nil {
			return nil, err
		}
	}

	switch node.Op {
	case ir.LogicalOpAnd:
		return NewValue(types.BOOL, val.Value.(bool) && right.Value.(bool)), nil

	case ir.LogicalOpOr:
		return NewValue(types.BOOL, val.Value.(bool) || right.Value.(bool)), nil

	case ir.LogicalOpNot:
		return NewValue(types.BOOL, !val.Value.(bool)), nil

	default:
		return nil, pos.Error("unknown logical operation: %d", node.Op)
	}
}
