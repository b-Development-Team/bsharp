package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
)

func (c *CGen) addFnCall(n *ir.FnCallNode) (*Code, error) {
	fn, err := c.AddNode(n.Fn)
	if err != nil {
		return nil, err
	}
	pre := fn.Pre

	call := &strings.Builder{}
	fmt.Fprintf(call, "%s(", fn.Value)
	for i, arg := range n.Args {
		if i != 0 {
			call.WriteString(", ")
		}
		arg, err := c.AddNode(arg)
		if err != nil {
			return nil, err
		}
		call.WriteString(arg.Value)
		pre = JoinCode(pre, fn.Pre)
	}
	call.WriteString(")")

	if isDynamic(n.Type()) { // Ret type is dynamic, free
		name := c.GetTmp("call")
		pre = JoinCode(pre, fmt.Sprintf("%s %s = %s;", c.CType(n.Type()), name, call.String()))
		c.stack.Add(c.FreeCode(name, n.Type()))

		return &Code{
			Pre:   pre,
			Value: name,
		}, nil
	}

	return &Code{
		Pre:   pre,
		Value: call.String(),
	}, nil
}

func (c *CGen) addFn(n *ir.FnNode) *Code {
	return &Code{
		Value: fmt.Sprintf("(&%s)", Namespace+n.Name),
	}
}

func (c *CGen) addReturn(n *ir.ReturnNode) (*Code, error) {
	v, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	c.isReturn = true
	// Grab if dynamic
	pre := JoinCode(c.stack.FreeCode(), v.Pre, fmt.Sprintf("return %s;", v.Value))
	if isDynamic(n.Value.Type()) {
		pre = JoinCode(c.GrabCode(v.Value, n.Value.Type()), pre)
	}
	return &Code{
		Pre: pre,
	}, nil
}
