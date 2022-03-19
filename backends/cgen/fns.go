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
	call.WriteString(");")

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
