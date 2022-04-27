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

	switch node.Lhs.Type().BasicType() {
	case types.INT:
		v := mathOp(left.Value.(int), right.Value.(int), node.Op)
		if v == nil {
			return nil, pos.Error("unknown math operation: %s", node.Op)
		}
		return NewValue(types.INT, *v), nil

	case types.FLOAT:
		v := mathOp(left.Value.(float64), right.Value.(float64), node.Op)
		if v == nil {
			return nil, pos.Error("unknown math operation: %s", node.Op)
		}
		return NewValue(types.INT, *v), nil

	default:
		return nil, pos.Error("cannot perform operations on type %s", left.Type.String())
	}
}

type mathOpValue interface {
	int | float64
}

func mathOp[T mathOpValue](lhs T, rhs T, op ir.MathOperation) *T {
	var out T
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
		switch any(lhs).(type) {
		case int:
			out = T(math.Pow(float64(lhs), float64(rhs)))

		case float64:
			out = T(math.Pow(float64(lhs), float64(rhs)))
		}

	case ir.MathOperationMod:
		switch any(lhs).(type) {
		case int:
			out = T(int(lhs) % int(rhs))

		case float64:
			out = T(math.Mod(float64(lhs), float64(rhs)))
		}

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
	switch node.Lhs.Type().BasicType() {
	case types.INT:
		v := compOp(left.Value.(int), right.Value.(int), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	case types.BYTE:
		v := compOp(left.Value.(byte), right.Value.(byte), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	case types.FLOAT:
		v := compOp(left.Value.(float64), right.Value.(float64), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	case types.STRING:
		v := compOp(left.Value.(string), right.Value.(string), node.Op)
		if v == nil {
			return nil, pos.Error("unknown compare operation: %s", node.Op)
		}
		return NewValue(types.BOOL, *v), nil

	default:
		return nil, pos.Error("cannot perform comparison on type %s", left.Type.String())
	}
}

type compOpValue interface {
	int | float64 | string | byte
}

func compOp[T compOpValue](lhs T, rhs T, op ir.CompareOperation) *bool {
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
