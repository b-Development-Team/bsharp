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
		out := &strings.Builder{}
		t := typ.(*types.FuncType)
		out.WriteString(c.CType(t.RetType))
		out.WriteString(" (*)(")
		for i, arg := range t.ParTypes {
			if i != 0 {
				out.WriteString(", ")
			}
			out.WriteString(c.CType(arg))
		}
		out.WriteString(")")
		return out.String()

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
