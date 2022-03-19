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
	}

	panic("invalid const")
}

func (c *CGen) FreeCode(varName string, typ types.Type) string {
	switch typ.BasicType() {
	case types.STRING:
		return fmt.Sprintf("string_free(%s);", varName)
	}

	panic("invalid type")
}

func (c *CGen) GrabCode(varName string, typ types.Type) string {
	switch typ.BasicType() {
	case types.STRING:
		return fmt.Sprintf("%s->refs++;", varName)
	}

	panic("invalid type")
}
