package cgen

import (
	"fmt"
	"strconv"

	"github.com/Nv7-Github/bsharp/ir"
)

func (c *CGen) addDefine(n *ir.DefineNode) (*Code, error) {
	v := c.ir.Variables[n.Var]
	name := Namespace + v.Name + strconv.Itoa(v.ID)
	val, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}

	pre := val.Pre
	free := ""
	if isDynamic(v.Type) { // Free old val, grab new one
		pre = JoinCode(pre, c.GrabCode(val.Value, n.Value.Type()))
		if c.declaredVars[v.ID] {
			oldName := c.GetTmp("old")
			pre = JoinCode(fmt.Sprintf("%s %s = %s;", c.CType(v.Type), oldName, name), pre)
			free = c.FreeCode(oldName, v.Type)
		}
	}

	code := fmt.Sprintf("%s = %s;", name, val.Value)
	// If not declared, then declare
	if !c.declaredVars[v.ID] {
		// Check if it is global
		if v.NeedsGlobal {
			fmt.Fprintf(c.globals, "%s %s;\n", c.CType(v.Type), name)
			c.declaredVars[v.ID] = true
		} else {
			code = fmt.Sprintf("%s %s = %s;", c.CType(v.Type), name, val.Value)
			c.declaredVars[v.ID] = true
		}

		if isDynamic(v.Type) {
			c.stack.Add(c.FreeCode(name, v.Type))
		}
	}

	return &Code{
		Pre: JoinCode(pre, code, free),
	}, nil
}
