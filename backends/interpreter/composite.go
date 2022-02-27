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

func (i *Interpreter) evalMake(pos *tokens.Pos, n *ir.MakeNode) (*Value, error) {
	typ := n.Type().(*types.MapType)
	var out interface{}
	switch typ.KeyType.BasicType() {
	case types.INT:
		out = make(map[int]interface{})

	case types.FLOAT:
		out = make(map[float64]interface{})

	case types.STRING:
		out = make(map[string]interface{})

	default:
		return nil, pos.Error("cannot make map with key type %s", typ.KeyType.String())
	}

	return NewValue(n.Type(), out), nil
}

func (i *Interpreter) evalSet(pos *tokens.Pos, n *ir.SetNode) (*Value, error) {
	m, err := i.evalNode(n.Map)
	if err != nil {
		return nil, err
	}
	k, err := i.evalNode(n.Key)
	if err != nil {
		return nil, err
	}
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	switch k.Type.BasicType() {
	case types.INT:
		m.Value.(map[int]interface{})[k.Value.(int)] = v.Value

	case types.FLOAT:
		m.Value.(map[float64]interface{})[k.Value.(float64)] = v.Value

	case types.STRING:
		m.Value.(map[string]interface{})[k.Value.(string)] = v.Value

	default:
		return nil, pos.Error("cannot set map with key type %s", k.Type.String())
	}

	return NewValue(types.NULL, nil), nil
}

func (i *Interpreter) evalGet(pos *tokens.Pos, n *ir.GetNode) (*Value, error) {
	m, err := i.evalNode(n.Map)
	if err != nil {
		return nil, err
	}
	k, err := i.evalNode(n.Key)
	if err != nil {
		return nil, err
	}
	var out interface{}
	switch k.Type.BasicType() {
	case types.INT:
		out = m.Value.(map[int]interface{})[k.Value.(int)]

	case types.FLOAT:
		out = m.Value.(map[float64]interface{})[k.Value.(float64)]

	case types.STRING:
		out = m.Value.(map[string]interface{})[k.Value.(string)]

	default:
		return nil, pos.Error("cannot get map with key type %s", k.Type.String())
	}

	return NewValue(n.Type(), out), nil
}
