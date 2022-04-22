package constrm

import (
	"math"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/old/ssa"
	"github.com/Nv7-Github/bsharp/types"
)

type compOpValue interface {
	int | float64 | string
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

func evalMath(blk *ssa.Block, i *ssa.Math) *ssa.Const {
	switch i.Type().BasicType() {
	case types.INT:
		return &ssa.Const{
			Value: mathOp(cnst(blk, i.Lhs).(int), cnst(blk, i.Rhs).(int), i.Op),
			Typ:   types.INT,
		}

	case types.FLOAT:
		return &ssa.Const{
			Value: mathOp(cnst(blk, i.Lhs).(float64), cnst(blk, i.Rhs).(float64), i.Op),
			Typ:   types.FLOAT,
		}
	}

	return nil
}

func evalComp(blk *ssa.Block, i *ssa.Compare) *ssa.Const {
	l := blk.Instructions[i.Lhs]
	var out *bool
	switch l.Type().BasicType() {
	case types.INT:
		out = compOp(cnst(blk, i.Lhs).(int), cnst(blk, i.Rhs).(int), i.Op)
	case types.FLOAT:
		out = compOp(cnst(blk, i.Lhs).(float64), cnst(blk, i.Rhs).(float64), i.Op)

	case types.STRING:
		out = compOp(cnst(blk, i.Lhs).(string), cnst(blk, i.Rhs).(string), i.Op)
	}
	if out != nil {
		return &ssa.Const{
			Value: *out,
			Typ:   types.BOOL,
		}
	}

	return nil
}
