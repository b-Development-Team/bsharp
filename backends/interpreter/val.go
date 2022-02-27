package interpreter

import (
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
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

	switch v.Type.BasicType() {
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
