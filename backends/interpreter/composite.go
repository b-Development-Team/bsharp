package interpreter

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func (i *Interpreter) evalIndex(pos *tokens.Pos, n *ir.IndexNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	indV, err := i.evalNode(n.Index)
	if err != nil {
		return nil, err
	}
	ind := indV.Value.(int)
	switch v.Type.BasicType() {
	case types.ARRAY:
		a := v.Value.([]interface{})
		if ind < 0 || ind >= len(a) {
			return nil, pos.Error("index out of bounds: %d with length %d", ind, len(a))
		}
		return NewValue(n.Type(), a[ind]), nil

	case types.STRING:
		a := []rune(v.Value.(string))
		if ind < 0 || ind >= len(a) {
			return nil, pos.Error("index out of bounds: %d with length %d", ind, len(a))
		}
		return NewValue(n.Type(), string(a[ind])), nil

	default:
		return nil, pos.Error("cannot index type %s", v.Type.String())
	}
}

func (i *Interpreter) evalLength(pos *tokens.Pos, n *ir.LengthNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	switch v.Type.BasicType() {
	case types.ARRAY:
		a := v.Value.([]interface{})
		return NewValue(types.INT, len(a)), nil

	case types.STRING:
		a := []rune(v.Value.(string))
		return NewValue(types.INT, len(a)), nil

	default:
		return nil, pos.Error("cannot get length of type %s", v.Type.String())
	}
}
