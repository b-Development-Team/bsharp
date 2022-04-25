package cgen

import (
	"fmt"
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

	case types.BOOL:
		return "bool"

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

	case types.STRUCT:
		name := typName(typ)
		_, exists := c.addedFns[name]
		if exists {
			return "struct " + name + "*"
		}
		sTyp := &strings.Builder{}
		fmt.Fprintf(sTyp, "struct %s {\n", name)
		for i, field := range typ.(*types.StructType).Fields {
			fmt.Fprintf(sTyp, "%s%s f%d;\n", c.Config.Tab, c.CType(field.Type), i)
		}
		fmt.Fprintf(sTyp, "%sint refs;\n", c.Config.Tab)
		sTyp.WriteString("};\n\n")
		c.globaltyps.WriteString(sTyp.String())

		// Free function
		out := &strings.Builder{}
		fmt.Fprintf(out, "void %s_free(struct %s* val) {;\n", name, name)
		fmt.Fprintf(out, "%sif (val == NULL) {\n%s%sreturn;\n%s}\n", c.Config.Tab, c.Config.Tab, c.Config.Tab, c.Config.Tab)
		fmt.Fprintf(out, "%sval->refs--;\n", c.Config.Tab)
		fmt.Fprintf(out, "%sif(val->refs == 0) {\n", c.Config.Tab)
		for i, field := range typ.(*types.StructType).Fields {
			if isDynamic(field.Type) {
				fmt.Fprintf(out, "%s%s%s;\n", c.Config.Tab, c.Config.Tab, c.FreeCode(fmt.Sprintf("val->f%d", i), field.Type))
			}
		}
		fmt.Fprintf(out, "%s%sfree(val);\n", c.Config.Tab, c.Config.Tab)
		fmt.Fprintf(out, "%s}\n", c.Config.Tab)
		out.WriteString("}\n\n")

		// Return
		c.addedFns[name] = struct{}{}
		c.globalfns.WriteString(out.String())
		return "struct " + name + "*"

	case types.ANY:
		return "any*"
	}

	panic("invalid type")
}
