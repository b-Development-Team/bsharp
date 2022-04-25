package cgen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (c *CGen) addGetStruct(n *ir.GetStructNode) (*Code, error) {
	s, err := c.AddNode(n.Struct)
	if err != nil {
		return nil, err
	}
	return &Code{
		Value: fmt.Sprintf("%s->f%d", s.Value, n.Field),
	}, nil
}

func (c *CGen) addSetStruct(n *ir.SetStructNode) (*Code, error) {
	s, err := c.AddNode(n.Struct)
	if err != nil {
		return nil, err
	}
	v, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}

	pre := v.Pre
	free := ""
	if isDynamic(n.Value.Type()) {
		// Need to grab value
		pre = JoinCode(pre, c.GrabCode(v.Value, n.Value.Type()))

		// Then, save old value, free after the new value is set (because new value might rely on old value)
		old := c.GetTmp("old")
		pre = JoinCode(s.Pre, fmt.Sprintf("%s %s = %s->f%d;", c.CType(n.Value.Type()), old, s.Value, n.Field), pre)
		free = c.FreeCode(old, n.Value.Type())
	} else {
		pre = JoinCode(s.Pre, pre)
	}

	// Set value
	code := fmt.Sprintf("%s->f%d = %s;", s.Value, n.Field, v.Value)
	return &Code{
		Pre: JoinCode(pre, code, free),
	}, nil
}

func (c *CGen) ZeroValue(t types.Type) string {
	switch t.BasicType() {
	case types.INT, types.FLOAT:
		return "0"

	case types.BOOL:
		return "false"

	case types.STRING, types.ARRAY, types.MAP, types.FUNCTION, types.STRUCT, types.ANY:
		return "NULL"
	}

	panic("invalid type") // might be NULL, but NULL should never be in a struct
}
