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
		a := *v.Value.(*[]any)
		if ind < 0 || ind >= len(a) {
			return nil, pos.Error("index out of bounds: %d with length %d", ind, len(a))
		}
		return NewValue(n.Type(), a[ind]), nil

	case types.STRING:
		a := v.Value.(string)
		if ind < 0 || ind >= len(a) {
			return nil, pos.Error("index out of bounds: %d with length %d", ind, len(a))
		}
		return NewValue(n.Type(), a[ind]), nil

	default:
		return nil, pos.Error("cannot index type %s", v.Type.String())
	}
}

func (i *Interpreter) evalSetIndex(pos *tokens.Pos, n *ir.SetIndexNode) (*Value, error) {
	a, err := i.evalNode(n.Array)
	if err != nil {
		return nil, err
	}
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	indV, err := i.evalNode(n.Index)
	if err != nil {
		return nil, err
	}
	ind := indV.Value.(int)
	if ind > len(*a.Value.(*[]any)) {
		return nil, pos.Error("index out of bounds: %d with length %d", ind, len(*a.Value.(*[]any)))
	}
	(*a.Value.(*[]any))[ind] = v.Value
	return NewValue(types.NULL, nil), nil
}

func (i *Interpreter) evalLength(pos *tokens.Pos, n *ir.LengthNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	switch v.Type.BasicType() {
	case types.ARRAY:
		a := *v.Value.(*[]any)
		return NewValue(types.INT, len(a)), nil

	case types.STRING:
		return NewValue(types.INT, len(v.Value.(string))), nil

	case types.MAP:
		switch v.Type.(*types.MapType).KeyType {
		case types.INT:
			return NewValue(types.INT, len(v.Value.(map[int]any))), nil

		case types.BYTE:
			return NewValue(types.INT, len(v.Value.(map[byte]any))), nil

		case types.FLOAT:
			return NewValue(types.INT, len(v.Value.(map[float64]any))), nil

		case types.STRING:
			return NewValue(types.INT, len(v.Value.(map[string]any))), nil
		}
	}

	return nil, pos.Error("cannot get length of type %s", v.Type.String())
}

func (i *Interpreter) evalMake(pos *tokens.Pos, n *ir.MakeNode) (*Value, error) {
	typ, ok := n.Type().(*types.MapType)
	if !ok { // Its an array if not map
		// Is struct
		if types.STRUCT.Equal(n.Type()) {
			t := n.Type().(*types.StructType)
			val := make([]any, len(t.Fields))
			return NewValue(n.Type(), &val), nil
		}

		v := make([]any, 0)
		return NewValue(n.Type(), &v), nil
	}
	var out any
	switch typ.KeyType.BasicType() {
	case types.INT:
		out = make(map[int]any)

	case types.BYTE:
		out = make(map[byte]any)

	case types.FLOAT:
		out = make(map[float64]any)

	case types.STRING:
		out = make(map[string]any)

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
		m.Value.(map[int]any)[k.Value.(int)] = v.Value

	case types.BYTE:
		m.Value.(map[byte]any)[k.Value.(byte)] = v.Value

	case types.FLOAT:
		m.Value.(map[float64]any)[k.Value.(float64)] = v.Value

	case types.STRING:
		m.Value.(map[string]any)[k.Value.(string)] = v.Value

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
	var out any
	var exists bool
	switch k.Type.BasicType() {
	case types.INT:
		out, exists = m.Value.(map[int]any)[k.Value.(int)]

	case types.BYTE:
		out, exists = m.Value.(map[byte]any)[k.Value.(byte)]

	case types.FLOAT:
		out, exists = m.Value.(map[float64]any)[k.Value.(float64)]

	case types.STRING:
		out, exists = m.Value.(map[string]any)[k.Value.(string)]

	default:
		return nil, pos.Error("cannot get map with key type %s", k.Type.String())
	}

	if !exists {
		return nil, pos.Error("key not in map: %v", k.Value)
	}

	return NewValue(n.Type(), out), nil
}

func (i *Interpreter) evalKeys(pos *tokens.Pos, n *ir.KeysNode) (*Value, error) {
	m, err := i.evalNode(n.Map)
	if err != nil {
		return nil, err
	}
	var out []any
	switch m.Type.(*types.MapType).KeyType.BasicType() {
	case types.INT:
		out = getKeys(m.Value.(map[int]any))

	case types.BYTE:
		out = getKeys(m.Value.(map[byte]any))

	case types.FLOAT:
		out = getKeys(m.Value.(map[float64]any))

	case types.STRING:
		out = getKeys(m.Value.(map[string]any))
	}

	return NewValue(n.Type(), &out), nil
}

func getKeys[T comparable](v map[T]any) []any {
	out := make([]any, len(v))
	i := 0
	for k := range v {
		out[i] = k
		i++
	}
	return out
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
		_, exists = m.Value.(map[int]any)[k.Value.(int)]

	case types.BYTE:
		_, exists = m.Value.(map[byte]any)[k.Value.(byte)]

	case types.FLOAT:
		_, exists = m.Value.(map[float64]any)[k.Value.(float64)]

	case types.STRING:
		_, exists = m.Value.(map[string]any)[k.Value.(string)]

	default:
		return nil, pos.Error("cannot get map with key type %s", k.Type.String())
	}

	return NewValue(types.BOOL, exists), nil
}

func (i *Interpreter) evalArray(n *ir.ArrayNode) (*Value, error) {
	out := make([]any, len(n.Values))
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
	v := arr.Value.(*[]any)

	val, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}

	*v = append(*v, val.Value)
	return NewValue(types.NULL, nil), nil
}

func (i *Interpreter) evalSlice(pos *tokens.Pos, n *ir.SliceNode) (*Value, error) {
	v, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	start, err := i.evalNode(n.Start)
	if err != nil {
		return nil, err
	}
	end, err := i.evalNode(n.End)
	if err != nil {
		return nil, err
	}
	if start.Value.(int) < 0 {
		return nil, pos.Error("start index out of bounds: %d", start.Value.(int))
	}
	if end.Value.(int) < start.Value.(int) {
		return nil, pos.Error("end index must be greater than start index")
	}

	switch v.Type.BasicType() {
	case types.ARRAY:
		a := v.Value.(*[]any)
		if start.Value.(int) >= len(*a) {
			return nil, pos.Error("start index out of bounds: %d with length %d", start.Value.(int), len(*a))
		}
		if end.Value.(int) > len(*a) {
			return nil, pos.Error("end index out of bounds: %d with length %d", end.Value.(int), len(*a))
		}
		*a = (*a)[start.Value.(int):end.Value.(int)]
		return NewValue(types.NULL, nil), nil

	case types.STRING:
		a := v.Value.(string)
		if start.Value.(int) >= len(a) {
			return nil, pos.Error("start index out of bounds: %d with length %d", start.Value.(int), len(a))
		}
		if end.Value.(int) > len(a) {
			return nil, pos.Error("end index out of bounds: %d with length %d", end.Value.(int), len(a))
		}
		return NewValue(types.STRING, a[start.Value.(int):end.Value.(int)]), nil
	}

	return nil, pos.Error("cannot get length of type %s", v.Type.String())
}

func (i *Interpreter) evalGetStruct(n *ir.GetStructNode) (*Value, error) {
	v, err := i.evalNode(n.Struct)
	if err != nil {
		return nil, err
	}
	return NewValue(n.Type(), (*v.Value.(*[]any))[n.Field]), nil
}

func (i *Interpreter) evalSetStruct(n *ir.SetStructNode) (*Value, error) {
	v, err := i.evalNode(n.Struct)
	if err != nil {
		return nil, err
	}
	val, err := i.evalNode(n.Value)
	if err != nil {
		return nil, err
	}
	(*v.Value.(*[]any))[n.Field] = val.Value
	return NewValue(types.NULL, nil), nil
}
