package cgen

import (
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
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

func (c *CGen) hashCase(v *ir.Const, typ types.Type) string {
	switch typ.BasicType() {
	case types.INT:
		return fmt.Sprintf("%d", v.Value.(int))

	case types.FLOAT:
		return fmt.Sprintf("%f", v.Value.(float64))

	case types.STRING:
		h := fnv.New32a()
		h.Write([]byte(v.Value.(string)))
		return fmt.Sprintf("%d", h.Sum32())

	default:
		panic("invalid hash type")
	}
}

func (c *CGen) addSwitch(i *ir.SwitchNode) (*Code, error) {
	cond, err := c.AddNode(i.Value)
	if err != nil {
		return nil, err
	}
	if types.STRING.Equal(i.Value.Type()) {
		cond.Value = fmt.Sprintf("str_hash(%s)", cond.Value)
	}

	bld := &strings.Builder{}
	fmt.Fprintf(bld, "switch (%s) {\n", cond.Value)
	for _, cs := range i.Cases {
		fmt.Fprintf(bld, "case %s:;\n", c.hashCase(cs.Value, i.Value.Type()))
		c.stack.Push()
		for _, node := range cs.Body {
			code, err := c.AddNode(node)
			if err != nil {
				return nil, err
			}
			c.addCode(bld, code)
		}
		c.addFree(bld)
		c.stack.Pop()
		bld.WriteString(c.Config.Tab + "break;\n")
	}

	if i.Default != nil {
		bld.WriteString("default:;\n")
		c.stack.Push()
		for _, node := range i.Default {
			code, err := c.AddNode(node)
			if err != nil {
				return nil, err
			}
			c.addCode(bld, code)
		}
		c.addFree(bld)
		c.stack.Pop()
	}
	bld.WriteString("}")

	return &Code{
		Pre: JoinCode(cond.Pre, bld.String()),
	}, nil
}
