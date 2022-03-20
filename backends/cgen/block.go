package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
)

func (c *CGen) addIf(i *ir.IfNode) (*Code, error) {
	cond, err := c.AddNode(i.Condition)
	if err != nil {
		return nil, err
	}
	bld := &strings.Builder{}
	fmt.Fprintf(bld, "if (%s) {\n", cond.Value)
	c.stack.Push()
	for _, node := range i.Body {
		code, err := c.AddNode(node)
		if err != nil {
			return nil, err
		}
		c.addCode(bld, code)
	}
	c.addFree(bld)
	c.stack.Pop()
	bld.WriteString("}")

	if i.Else != nil {
		bld.WriteString(" else {\n")
		c.stack.Push()
		for _, node := range i.Else {
			code, err := c.AddNode(node)
			if err != nil {
				return nil, err
			}
			c.addCode(bld, code)
		}
		c.addFree(bld)
		c.stack.Pop()
		bld.WriteString("}")
	}

	return &Code{
		Pre: JoinCode(cond.Pre, bld.String()),
	}, nil
}

func (c *CGen) addWhile(i *ir.WhileNode) (*Code, error) {
	cond, err := c.AddNode(i.Condition)
	if err != nil {
		return nil, err
	}
	bld := &strings.Builder{}
	fmt.Fprintf(bld, "while (%s) {\n", cond.Value)
	c.stack.Push()
	for _, node := range i.Body {
		code, err := c.AddNode(node)
		if err != nil {
			return nil, err
		}
		c.addCode(bld, code)
	}
	c.addFree(bld)
	c.stack.Pop()
	bld.WriteString("}")

	return &Code{
		Pre: JoinCode(cond.Pre, bld.String()),
	}, nil
}
