package cgen

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

var mathOps = map[ir.MathOperation]string{
	ir.MathOperationAdd: "+",
	ir.MathOperationSub: "-",
	ir.MathOperationMul: "*",
	ir.MathOperationDiv: "/",
	ir.MathOperationMod: "%",
}

func (c *CGen) addMath(n *ir.MathNode) (*Code, error) {
	l, err := c.AddNode(n.Lhs)
	if err != nil {
		return nil, err
	}
	r, err := c.AddNode(n.Rhs)
	if err != nil {
		return nil, err
	}
	pre := JoinCode(l.Pre, r.Pre)
	if n.Op == ir.MathOperationPow {
		switch n.Lhs.Type().BasicType() {
		case types.INT:
			return &Code{
				Pre:   pre,
				Value: fmt.Sprintf("(long)(pow((double)%s, (double)%s) + 0.5)", l.Value, r.Value),
			}, nil

		case types.FLOAT:
			return &Code{
				Pre:   pre,
				Value: fmt.Sprintf("pow(%s, %s)", l.Value, r.Value),
			}, nil
		}
	}
	if n.Op == ir.MathOperationMod && types.FLOAT.Equal(n.Lhs.Type()) { // Need to use fmod
		return &Code{
			Pre:   pre,
			Value: fmt.Sprintf("fmod(%s, %s)", l.Value, r.Value),
		}, nil
	}
	return &Code{
		Pre:   pre,
		Value: fmt.Sprintf("(%s %s %s)", l.Value, mathOps[n.Op], r.Value),
	}, nil
}

var compOpCodes = map[ir.CompareOperation]string{
	ir.CompareOperationEqual:        "==",
	ir.CompareOperationNotEqual:     "!=",
	ir.CompareOperationGreater:      ">",
	ir.CompareOperationGreaterEqual: ">=",
	ir.CompareOperationLess:         "<",
	ir.CompareOperationLessEqual:    "<=",
}

func (c *CGen) addCompare(n *ir.CompareNode) (*Code, error) {
	l, err := c.AddNode(n.Lhs)
	if err != nil {
		return nil, err
	}
	r, err := c.AddNode(n.Rhs)
	if err != nil {
		return nil, err
	}

	switch n.Lhs.Type().BasicType() {
	case types.INT, types.FLOAT, types.BYTE:
		return &Code{
			Pre:   JoinCode(l.Pre, r.Pre),
			Value: fmt.Sprintf("(%s %s %s)", l.Value, compOpCodes[n.Op], r.Value),
		}, nil

	case types.STRING:
		return &Code{
			Pre:   JoinCode(l.Pre, r.Pre),
			Value: fmt.Sprintf("string_cmp(%s, %s) %s 0", l.Value, r.Value, compOpCodes[n.Op]),
		}, nil
	}

	panic("invalid type")
}
