package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

func typName(typ types.Type) string {
	switch typ.BasicType() {
	case types.INT:
		return "i"

	case types.FLOAT:
		return "f"

	case types.STRING:
		return "s"

	case types.BOOL:
		return "b"

	case types.ARRAY:
		return "a" + typName(typ.(*types.ArrayType).ElemType)

	case types.MAP:
		return "m" + typName(typ.(*types.MapType).KeyType) + typName(typ.(*types.MapType).ValType)

	case types.FUNCTION:
		out := &strings.Builder{}
		t := typ.(*types.FuncType)
		out.WriteString("f")
		for _, par := range t.ParTypes {
			out.WriteString(typName(par))
		}
		out.WriteString(typName(t.RetType))
		return out.String()

	case types.NULL:
		return "n"

	case types.STRUCT:
		out := &strings.Builder{}
		t := typ.(*types.StructType)
		out.WriteString("s")
		for _, f := range t.Fields {
			out.WriteString(typName(f.Type))
		}
		return out.String()

	default:
		panic("unknown")
	}
}

func (c *CGen) arrFreeFn(typ types.Type) string {
	e := typ.(*types.ArrayType).ElemType
	if !isDynamic(e) {
		return "NULL"
	}

	name := "arrfree_" + typName(typ)
	_, exists := c.addedFns[name]
	if exists {
		return "&" + name
	}

	c.addedFns[name] = struct{}{}

	code := &strings.Builder{}
	fmt.Fprintf(code, "void %s(array* arr) {\n", name)
	fmt.Fprintf(code, "%sfor (int i = 0; i < arr->len; i++) {\n", c.Config.Tab)
	code.WriteString(c.Config.Tab + c.Config.Tab + c.FreeCode("*(("+c.CType(e)+"*)(array_get(arr, i)))", e) + "\n")
	fmt.Fprintf(code, "%s}\n}\n\n", c.Config.Tab)

	c.globalfns.WriteString(code.String())

	return "&" + name
}

func (c *CGen) addArray(n *ir.ArrayNode) (*Code, error) {
	arr := c.GetTmp("arr")
	pre := fmt.Sprintf("array* %s = array_new(sizeof(%s), %d);", arr, c.CType(n.Type().(*types.ArrayType).ElemType), len(n.Values))
	for _, v := range n.Values {
		val, err := c.AddNode(v)
		if err != nil {
			return nil, err
		}
		pre = JoinCode(pre, val.Pre)
		if isDynamic(v.Type()) {
			pre = JoinCode(pre, c.GrabCode(val.Value, v.Type()))
		} else { // Need to be able to get pointer
			name := c.GetTmp("cnst")
			pre = JoinCode(pre, fmt.Sprintf("%s %s = %s;", c.CType(v.Type()), name, val.Value))
			val.Value = name
		}
		pre = JoinCode(pre, fmt.Sprintf("array_append(%s, &(%s));", arr, val.Value))
	}
	c.stack.Add(c.FreeCode(arr, n.Type()))
	return &Code{Pre: pre, Value: arr}, nil
}

func (c *CGen) addIndex(pos *tokens.Pos, n *ir.IndexNode) (*Code, error) {
	arr, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	ind, err := c.AddNode(n.Index)
	if err != nil {
		return nil, err
	}
	if types.STRING.Equal(n.Value.Type()) {
		name := c.GetTmp("ind")
		pre := fmt.Sprintf("string* %s = string_ind(%s, %s);", name, arr.Value, ind.Value)
		c.stack.Add(c.FreeCode(name, types.STRING))
		// Bounds check
		if c.Config.BoundsCheck {
			pre = JoinCode(fmt.Sprintf("string_bounds(%s, %s, %q);", arr.Value, ind.Value, pos.String()), pre)
		}
		return &Code{
			Pre:   JoinCode(arr.Pre, ind.Pre, pre),
			Value: name,
		}, nil
	}
	typ := c.CType(n.Type())
	pre := JoinCode(arr.Pre, ind.Pre)
	// Bounds check
	if c.Config.BoundsCheck {
		pre = JoinCode(pre, fmt.Sprintf("array_bounds(%s, %s, %q);", arr.Value, ind.Value, pos.String()))
	}
	return &Code{
		Pre:   pre,
		Value: fmt.Sprintf("(*((%s*)(array_get(%s, %s))))", typ, arr.Value, ind.Value),
	}, nil
}

func (c *CGen) addAppend(n *ir.AppendNode) (*Code, error) {
	arr, err := c.AddNode(n.Array)
	if err != nil {
		return nil, err
	}
	val, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	pre := JoinCode(arr.Pre, val.Pre)
	if isDynamic(n.Array.Type().(*types.ArrayType).ElemType) {
		pre = JoinCode(pre, c.GrabCode(val.Value, n.Value.Type()))
	}

	// Need to be able to get pointer
	if !isDynamic(n.Value.Type()) {
		name := c.GetTmp("cnst")
		pre = JoinCode(pre, fmt.Sprintf("%s %s = %s;", c.CType(n.Value.Type()), name, val.Value))
		val.Value = name
	}
	pre = JoinCode(pre, fmt.Sprintf("array_append(%s, &(%s));", arr.Value, val.Value))
	return &Code{
		Pre: pre,
	}, nil
}

func (c *CGen) addLength(n *ir.LengthNode) (*Code, error) {
	arr, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	// Works for every type
	return &Code{
		Pre:   JoinCode(arr.Pre),
		Value: fmt.Sprintf("%s->len", arr.Value),
	}, nil
}

func (c *CGen) addSlice(n *ir.SliceNode) (*Code, error) {
	v, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	start, err := c.AddNode(n.Start)
	if err != nil {
		return nil, err
	}
	end, err := c.AddNode(n.End)
	if err != nil {
		return nil, err
	}
	if types.STRING.Equal(n.Value.Type()) {
		name := c.GetTmp("slice")
		pre := fmt.Sprintf("string* %s = string_slice(%s, %s, %s);", name, v.Value, start.Value, end.Value)
		c.stack.Add(c.FreeCode(name, types.STRING))
		return &Code{
			Pre:   JoinCode(v.Pre, start.Pre, end.Pre, pre),
			Value: name,
		}, nil
	}

	// Array
	code := fmt.Sprintf("array_slice(%s, %s, %s);", v.Value, start.Value, end.Value)
	return &Code{
		Pre:   JoinCode(v.Pre, start.Pre, end.Pre),
		Value: code,
	}, nil
}

func (c *CGen) addSetIndex(pos *tokens.Pos, n *ir.SetIndexNode) (*Code, error) {
	arr, err := c.AddNode(n.Array)
	if err != nil {
		return nil, err
	}
	ind, err := c.AddNode(n.Index)
	if err != nil {
		return nil, err
	}
	val, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}

	pre := val.Pre
	check := ""
	if c.Config.BoundsCheck {
		check = fmt.Sprintf("array_bounds(%s, %s, %q);", arr.Value, ind.Value, pos.String())
	}
	free := ""
	if isDynamic(n.Value.Type()) {
		pre = JoinCode(pre, c.GrabCode(val.Value, n.Value.Type()))
		old := c.GetTmp("old")
		pre = JoinCode(arr.Pre, ind.Pre, check, fmt.Sprintf("%s %s = (*((%s*)(array_get(%s, %s))));", c.CType(n.Value.Type()), old, c.CType(n.Value.Type()), arr.Value, ind.Value), pre)
		free = c.FreeCode(old, n.Value.Type())
	} else {
		pre = JoinCode(arr.Pre, ind.Pre, check, pre)
	}

	code := fmt.Sprintf("array_set(%s, %s, &(%s));", arr.Value, ind.Value, val.Value)
	return &Code{
		Pre: JoinCode(pre, code, free),
	}, nil
}
