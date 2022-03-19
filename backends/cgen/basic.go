package cgen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
)

func (c *CGen) addPrint(n *ir.PrintNode) (*Code, error) {
	v, err := c.AddNode(n.Arg)
	if err != nil {
		return nil, err
	}
	return &Code{
		Pre: JoinCode(v.Pre, fmt.Sprintf("string_print(%s);", v.Value)),
	}, nil
}
