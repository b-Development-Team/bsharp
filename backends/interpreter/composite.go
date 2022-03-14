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
		a := *v.Value.(*[]interface{})
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
		a := *v.Value.(*[]interface{})
		return NewValue(types.INT, len(a)), nil

	case types.STRING:
		a := []rune(v.Value.(string))
		return NewValue(types.INT, len(a)), nil

	default:
		return nil, pos.Error("cannot get length of type %s", v.Type.String())
	}
}

func (i *Interpreter) evalMake(pos *tokens.Pos, n *ir.MakeNode) (*Value, error) {
	typ, ok := n.Type().(*types.MapType)
	if !ok { // Its an array if not map
		v := make([]interface{}, 0)
		return NewValue(n.Type(), &v), nil
	}
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

func (i *Interpreter) evalKeys(pos *tokens.Pos, n *ir.KeysNode) (*Value, error) {
	m, err := i.evalNode(n.Map)
	if err != nil {
		return nil, err
	}
	var out []interface{}
	switch m.Type.(*types.MapType).KeyType.BasicType() {
	case types.INT:
		v := m.Value.(map[int]interface{})
		out = make([]interface{}, len(v))
		i := 0
		for k := range v {
			out[i] = k
			i++
		}

	case types.FLOAT:
		v := m.Value.(map[float64]interface{})
		out = make([]interface{}, len(v))
		i := 0
		for k := range v {
			out[i] = k
			i++
		}

	case types.STRING:
		v := m.Value.(map[string]interface{})
		out = make([]interface{}, len(v))
		i := 0
		for k := range v {
			out[i] = k
			i++
		}
	}

	return NewValue(n.Type(), &out), nil
}

func (i *Interpreter) evalExists(pos *tokens.Pos, n *ir.ExistsNode) (*Value, error) {
	m, err := i.evalNode(n.Map)
	if err != nil {
		return nil, err
	}
	k, err := i.evalNode(n.Key)
	if err != nil {
		return nil, err
	}
	var exists bool
	switch k.Type.BasicType() {
	case types.INT:
		_, exists = m.Value.(map[int]interface{})[k.Value.(int)]

	case types.FLOAT:
		_, exists = m.Value.(map[float64]interface{})[k.Value.(float64)]

	case types.STRING:
		_, exists = m.Value.(map[string]interface{})[k.Value.(string)]

	default:
		return nil, pos.Error("cannot get map with key type %s", k.Type.String())
	}

	return NewValue(types.BOOL, exists), nil
}

func (i *Interpreter) evalArray(n *ir.ArrayNode) (*Value, error) {
	out := make([]interface{}, len(n.Values))
	for ind, v := range n.Values {
		val, err := i.evalNode(v)
		if err != nil {
			return nil, err
		}
		out[ind] = val.Value
	}
	return NewValue(types.ARRAY, &out), nil
}

func (i *Interpreter) evalAppend(n *ir.AppendNode) (*Value, error) {
	arr, err := i.evalNode(n.Array)
	if err != nil {
		return nil, err
	}
	v := arr.Value.(*[]interface{})

	val, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}

	*v = append(*v, val.Value)
	return NewValue(types.NULL, nil), nil
}
