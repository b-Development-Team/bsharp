package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
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

func (c *CGen) addPanic(n *ir.PanicNode, pos *tokens.Pos) (*Code, error) {
	v, err := c.AddNode(n.Arg)
	if err != nil {
		return nil, err
	}
	return &Code{
		Pre: JoinCode(v.Pre, fmt.Sprintf("bsp_panic(%s, %q);", v.Value, pos.String())),
	}, nil
}

var timeModeNames = map[ir.TimeMode]string{
	ir.TimeModeSeconds: "normaltime",
	ir.TimeModeMilli:   "millitime",
	ir.TimeModeMicro:   "microtime",
	ir.TimeModeNano:    "nanotime",
}

func (c *CGen) addTime(n *ir.TimeNode) *Code {
	return &Code{
		Value: fmt.Sprintf("%s()", timeModeNames[n.Mode]),
	}
}

func (c *CGen) addConcat(n *ir.ConcatNode) (*Code, error) {
	out := &strings.Builder{}
	fmt.Fprintf(out, "string_concat(%d", len(n.Values))
	pre := ""
	for _, v := range n.Values {
		v, err := c.AddNode(v)
		if err != nil {
			return nil, err
		}
		pre = JoinCode(pre, v.Pre)
		fmt.Fprintf(out, ", %s", v.Value)
	}

	name := c.GetTmp("concat")
	pre = JoinCode(pre, fmt.Sprintf("string* %s = %s);", name, out.String()))
	c.stack.Add(c.FreeCode(name, types.STRING))
	return &Code{
		Pre:   pre,
		Value: name,
	}, nil
}

func (c *CGen) addLogicalOp(n *ir.LogicalOpNode) (*Code, error) {
	l, err := c.AddNode(n.Val)
	if err != nil {
		return nil, err
	}
	var r *Code
	if n.Rhs != nil {
		r, err = c.AddNode(n.Rhs)
		if err != nil {
			return nil, err
		}
	}

	switch n.Op {
	case ir.LogicalOpAnd:
		return &Code{
			Pre:   JoinCode(l.Pre, r.Pre),
			Value: fmt.Sprintf("(%s && %s)", l.Value, r.Value),
		}, nil

	case ir.LogicalOpOr:
		return &Code{
			Pre:   JoinCode(l.Pre, r.Pre),
			Value: fmt.Sprintf("(%s || %s)", l.Value, r.Value),
		}, nil

	case ir.LogicalOpNot:
		return &Code{
			Pre:   JoinCode(l.Pre),
			Value: fmt.Sprintf("(!%s)", l.Value),
		}, nil
	}

	panic("invalid op")
}
