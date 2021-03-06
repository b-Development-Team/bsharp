package cgen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (c *CGen) addConst(n *ir.Const) *Code {
	switch n.Type().BasicType() {
	case types.INT:
		return &Code{
			Value: fmt.Sprintf("%d", n.Value),
		}

	case types.BYTE:
		return &Code{
			Value: fmt.Sprintf("%d", n.Value),
		}

	case types.FLOAT:
		return &Code{
			Value: fmt.Sprintf("%f", n.Value),
		}

	case types.STRING:
		name := c.GetTmp("string")
		c.stack.Add(c.FreeCode(name, types.STRING))
		return &Code{
			Pre:   fmt.Sprintf("string* %s = string_from_const(%q);", name, n.Value),
			Value: name,
		}

	case types.BOOL:
		if n.Value.(bool) {
			return &Code{
				Value: "true",
			}
		}
		return &Code{
			Value: "false",
		}

	case types.NULL:
		return &Code{
			Value: "NULL",
		}
	}

	panic("invalid const")
}

func (c *CGen) FreeCode(varName string, typ types.Type) string {
	switch typ.BasicType() {
	case types.STRING:
		return fmt.Sprintf("string_free(%s);", varName)

	case types.ARRAY:
		return fmt.Sprintf("array_free(%s, %s);", varName, c.arrFreeFn(typ))

	case types.MAP:
		return fmt.Sprintf("map_free(%s);", varName)

	case types.STRUCT:
		return fmt.Sprintf("%s_free(%s);", typName(typ), varName)

	case types.ANY:
		return fmt.Sprintf("any_free(%s);", varName)
	}

	panic("invalid type")
}

func (c *CGen) GrabCode(varName string, typ types.Type) string {
	switch typ.BasicType() {
	case types.STRING, types.ARRAY, types.MAP, types.STRUCT, types.ANY:
		return fmt.Sprintf("%s->refs++;", varName)
	}

	panic("invalid type")
}

func (c *CGen) addCast(n *ir.CastNode) (*Code, error) {
	if types.ANY.Equal(n.Type()) {
		return c.addAnyCast(n)
	}

	if types.ANY.Equal(n.Value.Type()) {
		return c.addAnyFromCast(n)
	}

	v, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}

	switch n.Value.Type().BasicType() {
	case types.INT:
		switch n.Type().BasicType() {
		case types.FLOAT:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("(double)(%s)", v.Value),
			}, nil

		case types.BYTE:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("(char)(%s)", v.Value),
			}, nil

		case types.STRING:
			name := c.GetTmp("itoa")
			pre := JoinCode(v.Pre, fmt.Sprintf("string* %s = string_itoa(%s);", name, v.Value))
			c.stack.Add(c.FreeCode(name, types.STRING))
			return &Code{
				Pre:   pre,
				Value: name,
			}, nil
		}

	case types.BYTE:
		switch n.Type().BasicType() {
		case types.INT:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("(long)(%s)", v.Value),
			}, nil

		case types.STRING:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("byte_to_string(%s)", v.Value),
			}, nil
		}

	case types.ARRAY: // ARRAY{BYTE}
		return &Code{
			Pre:   v.Pre,
			Value: fmt.Sprintf("array_byte_to_string(%s)", v.Value),
		}, nil

	case types.FLOAT:
		switch n.Type().BasicType() {
		case types.INT:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("(long)(%s)", v.Value),
			}, nil

		case types.STRING:
			name := c.GetTmp("ftoa")
			pre := JoinCode(v.Pre, fmt.Sprintf("string* %s = string_ftoa(%s);", name, v.Value))
			c.stack.Add(c.FreeCode(name, types.STRING))
			return &Code{
				Pre:   pre,
				Value: name,
			}, nil
		}

	case types.STRING:
		switch n.Type().BasicType() {
		case types.INT:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("bsharp_atoi(%s)", v.Value),
			}, nil

		case types.FLOAT:
			return &Code{
				Pre:   v.Pre,
				Value: fmt.Sprintf("bsharp_atof(%s)", v.Value),
			}, nil
		}

	case types.BOOL:
		switch n.Type().BasicType() {
		case types.STRING:
			name := c.GetTmp("btoa")
			pre := JoinCode(v.Pre, fmt.Sprintf("string* %s = string_btoa(%s);", name, v.Value))
			c.stack.Add(c.FreeCode(name, types.STRING))
			return &Code{
				Pre:   pre,
				Value: name,
			}, nil
		}
	}

	panic("cannot cast")
}

func isDynamic(typ types.Type) bool {
	switch typ.BasicType() {
	case types.STRING, types.MAP, types.ARRAY, types.STRUCT, types.ANY:
		return true
	}
	return false
}
