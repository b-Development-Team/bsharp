package interpreter

import (
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type anyVal struct {
	typ types.Type
	val any
}

func (i *Interpreter) evalConst(c *ir.Const) (*Value, error) {
	return NewValue(c.Type(), c.Value), nil
}

func (i *Interpreter) canCastAny(v interface{}, t types.Type) bool {
	return v.(anyVal).typ.Equal(t)
}

func (i *Interpreter) evalCanCast(pos *tokens.Pos, c *ir.CanCastNode) (*Value, error) {
	v, err := i.evalNode(c.Value)
	if err != nil {
		return nil, err
	}
	return NewValue(types.BOOL, i.canCastAny(v.Value, c.Typ)), nil
}

func (i *Interpreter) evalCast(c *ir.CastNode) (*Value, error) {
	v, err := i.evalNode(c.Value)
	if err != nil {
		return nil, err
	}

	if types.ANY.Equal(c.Type()) {
		return NewValue(types.ANY, anyVal{typ: v.Type, val: v.Value}), nil
	}

	if types.ANY.Equal(c.Value.Type()) {
		if !i.canCastAny(v.Value, c.Type()) {
			return nil, c.Pos().Error("cannot cast value of type %s to %s", v.Value.(anyVal).typ.String(), c.Type().String())
		}

		return NewValue(c.Type(), v.Value.(anyVal).val), nil
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

	case types.BYTE:
		switch c.Type().BasicType() {
		case types.INT:
		case types.BYTE:
			return NewValue(c.Type(), int(v.Value.(byte))), nil
		}
		fallthrough

	case types.ARRAY: // ARRAY{BYTE}
		v := *v.Value.(*[]any)
		vals := make([]byte, len(v))
		for i, val := range v {
			vals[i] = val.(byte)
		}
		return NewValue(c.Type(), string(vals)), nil

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
				return nil, c.Pos().Error("cannot cast %q to INT", v.Value.(string))
			}
			return NewValue(c.Type(), val), nil

		case types.FLOAT:
			val, err := strconv.ParseFloat(v.Value.(string), 64)
			if err != nil {
				return nil, c.Pos().Error("cannot cast %q to FLOAT", v.Value.(string))
			}
			return NewValue(c.Type(), val), nil

		case types.STRING:
			return NewValue(c.Type(), string(v.Value.(byte))), nil
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
