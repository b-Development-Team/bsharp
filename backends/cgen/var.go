package cgen

import (
	"fmt"
	"strconv"

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
		c.declaredVars[v.ID] = true
		code = fmt.Sprintf("%s %s = %s;", c.CType(v.Type), name, val.Value)
		c.stack.Add(c.FreeCode(name, v.Type))
	}

	return &Code{
		Pre: JoinCode(pre, code),
	}, nil
}
