package interpreter

import (
	"math"
	"math/rand"
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalConst(c *ir.Const) (*Value, error) {
	return NewValue(c.Type(), c.Value), nil
}

func (i *Interpreter) evalCast(c *ir.CastNode) (*Value, error) {
	v, err := i.evalNode(c.Value)
	if err != nil {
		return nil, err
	}

	switch c.Value.Type().BasicType() {
	case types.INT:
		switch c.Type().BasicType() {
		case types.FLOAT:
			return NewValue(c.Type(), float64(v.Value.(int))), nil

		case types.STRING:
			return NewValue(c.Type(), strconv.Itoa(v.Value.(int))), nil
		}
		fallthrough

	case types.FLOAT:
		switch c.Type().BasicType() {
		case types.INT:
			return NewValue(c.Type(), int(v.Value.(float64))), nil

		case types.STRING:
			return NewValue(c.Type(), strconv.FormatFloat(v.Value.(float64), 'f', -1, 64)), nil
		}
		fallthrough

	case types.STRING:
		switch c.Type().BasicType() {
		case types.INT:
			val, err := strconv.Atoi(v.Value.(string))
			if err != nil {
				return nil, c.Pos().Error("cannot cast \"%s\" to INT", v.Value.(string))
			}
			return NewValue(c.Type(), val), nil

		case types.FLOAT:
			val, err := strconv.ParseFloat(v.Value.(string), 64)
			if err != nil {
				return nil, c.Pos().Error("cannot cast \"%s\" to FLOAT", v.Value.(string))
			}
			return NewValue(c.Type(), val), nil
		}
		fallthrough

	case types.BOOL:
		switch c.Type().BasicType() {
		case types.STRING:
			val := "false"
			if v.Value.(bool) {
				val = "true"
			}
			return NewValue(c.Type(), val), nil
		}
		fallthrough

	default:
		return nil, c.Pos().Error("cannot cast type \"%s\" to \"%s\"", v.Type.BasicType(), c.Type().BasicType())
	}
}

func (i *Interpreter) evalRandint(n *ir.RandintNode) (*Value, error) {
	lower, err := i.evalNode(n.Lower)
	if err != nil {
		return nil, err
	}
	upper, err := i.evalNode(n.Upper)
	if err != nil {
		return nil, err
	}
	min := lower.Value.(int)
	max := upper.Value.(int)
	return NewValue(types.INT, rand.Intn(max-min+1)+min), nil
}

func (i *Interpreter) evalRandom(n *ir.RandomNode) (*Value, error) {
	lower, err := i.evalNode(n.Lower)
	if err != nil {
		return nil, err
	}
	upper, err := i.evalNode(n.Upper)
	if err != nil {
		return nil, err
	}
	min := lower.Value.(float64)
	max := upper.Value.(float64)
	return NewValue(types.FLOAT, min+rand.Float64()*(max-min)), nil
}

func (i *Interpreter) evalMathFn(pos *tokens.Pos, n *ir.MathFunctionNode) (*Value, error) {
	arg, err := i.evalNode(n.Arg)
	if err != nil {
		return nil, err
	}
	switch n.Func {
	case ir.MathFunctionCeil:
		return NewValue(types.INT, int(math.Ceil(arg.Value.(float64)))), nil

	case ir.MathFunctionFloor:
		return NewValue(types.INT, int(math.Floor(arg.Value.(float64)))), nil

	case ir.MathFunctionRound:
		return NewValue(types.INT, int(math.Round(arg.Value.(float64)))), nil

	default:
		return nil, pos.Error("unknown math function \"%s\"", n.Func)
	}
}
