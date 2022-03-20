package cgen

import (
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

const Namespace = "bsharp__"

func (c *CGen) CType(typ types.Type) string {
	switch typ.BasicType() {
	case types.INT:
		return "long"

	case types.FLOAT:
		return "double"

	case types.STRING:
		return "string*"

	case types.FUNCTION:
		name := typName(typ) + "_typ"
		_, exists := c.addedFns[name]
		if exists {
			return name
		}

		t := typ.(*types.FuncType)
		out := &strings.Builder{}
		out.WriteString("typedef ")
		out.WriteString(c.CType(t.RetType))
		out.WriteString(" (*")
		out.WriteString(name)
		out.WriteString(")(")
		for i, arg := range t.ParTypes {
			if i != 0 {
				out.WriteString(", ")
			}
			out.WriteString(c.CType(arg))
		}
		out.WriteString(");\n")
		c.globaltyps.WriteString(out.String())
		c.addedFns[name] = struct{}{}

		return name

	case types.ARRAY:
		return "array*"

	case types.MAP:
		return "map*"

	case types.NULL:
		return "void"

	default:
		return "unknown"
	}
}
