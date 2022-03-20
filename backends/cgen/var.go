package cgen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func isDynamic(typ types.Type) bool {
	switch typ.BasicType() {
	case types.STRING, types.MAP, types.ARRAY:
		return true
	}
	return false
}

func (c *CGen) fnTypeName(t *types.FuncType, name string) string {
	out := &strings.Builder{}
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
	out.WriteString(")")
	return out.String()
}

func (c *CGen) addDefine(n *ir.DefineNode) (*Code, error) {
	v := c.ir.Variables[n.Var]
	name := Namespace + v.Name + strconv.Itoa(v.ID)
	val, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}

	pre := val.Pre
	if isDynamic(v.Type) { // Free old val, grab new one
		pre = JoinCode(pre, c.GrabCode(val.Value, n.Value.Type()))
		if c.declaredVars[v.ID] {
			pre = JoinCode(c.FreeCode(name, v.Type), pre)
		}
	}

	code := fmt.Sprintf("%s = %s;", name, val.Value)
	// If not declared, then declare
	if !c.declaredVars[v.ID] {
		// Check if it is global
		if v.ScopeType == ir.ScopeTypeGlobal {
			if types.FUNCTION.Equal(v.Type) {
				t := v.Type.(*types.FuncType)
				fmt.Fprintf(c.globals, "%s;\n", c.fnTypeName(t, name))
			} else {
				fmt.Fprintf(c.globals, "%s %s;\n", c.CType(v.Type), name)
			}
			c.declaredVars[v.ID] = true
		} else {
			if types.FUNCTION.Equal(v.Type) {
				co := &strings.Builder{}
				t := v.Type.(*types.FuncType)
				co.WriteString(c.fnTypeName(t, name))
				co.WriteString(" = ")
				co.WriteString(val.Value)
				co.WriteString(";")
				code = co.String()
			} else {
				code = fmt.Sprintf("%s %s = %s;", c.CType(v.Type), name, val.Value)
			}
			c.declaredVars[v.ID] = true
		}

		if isDynamic(v.Type) {
			c.stack.Add(c.FreeCode(name, v.Type))
		}
	}

	return &Code{
		Pre: JoinCode(pre, code),
	}, nil
}
