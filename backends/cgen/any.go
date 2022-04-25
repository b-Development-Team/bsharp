package cgen

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

func (c *CGen) typID(typ types.Type) int {
	name := typName(typ)
	id, exists := c.typIDs[name]
	if !exists {
		id = len(c.typIDs)
		c.typIDs[name] = id
		c.idTyps = append(c.idTyps, typ)
	}
	return id
}

func (c *CGen) typCode() string {
	out := &strings.Builder{}
	out.WriteString("const char* const anytyps[] = {\n")
	for _, t := range c.idTyps {
		fmt.Fprintf(out, "%s%q,\n", c.Config.Tab, t.String())
	}
	out.WriteString("};")
	return out.String()
}

func (c *CGen) addAnyCast(n *ir.CastNode) (*Code, error) {
	v, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}

	freefn := "NULL"
	pre := v.Pre
	if isDynamic(n.Value.Type()) {
		// Grab value
		pre = JoinCode(pre, c.GrabCode(v.Value, n.Value.Type()))

		// Create free fn if not added
		name := typName(n.Value.Type()) + "_anyfree"
		_, exists := c.addedFns[name]
		if !exists {
			fmt.Fprintf(c.globalfns, "void %s(void* val) {\n", name)
			fmt.Fprintf(c.globalfns, "%s%s v = *((%s*)val);\n", c.Config.Tab, c.CType(n.Value.Type()), c.CType(n.Value.Type()))
			fmt.Fprintf(c.globalfns, "%s%s;\n", c.Config.Tab, c.FreeCode("v", n.Value.Type()))
			c.globalfns.WriteString("}\n\n")
			c.addedFns[name] = struct{}{}
		}
		freefn = name
	} else {
		// Need to be able to get pointer
		name := c.GetTmp("cnst")
		pre = JoinCode(pre, fmt.Sprintf("%s %s = %s;", c.CType(n.Value.Type()), name, v.Value))
		v.Value = name
	}

	a := c.GetTmp("any")
	code := fmt.Sprintf("any* %s = any_new((void*)(&%s), %d, %s);\n", a, v.Value, c.typID(n.Value.Type()), freefn)
	c.stack.Add(c.FreeCode(a, types.ANY))
	return &Code{
		Pre:   JoinCode(pre, code),
		Value: a,
	}, nil
}

func (c *CGen) addCanCast(n *ir.CanCastNode) (*Code, error) {
	a, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	t := c.typID(n.Typ)
	return &Code{
		Pre:   a.Pre,
		Value: fmt.Sprintf("(%s->typ == %d)", a.Value, t),
	}, nil
}

func (c *CGen) addAnyFromCast(n *ir.CastNode) (*Code, error) {
	a, err := c.AddNode(n.Value)
	if err != nil {
		return nil, err
	}
	pre := JoinCode(a.Pre, fmt.Sprintf("any_try_cast(%q, %s, %d);", n.Pos().String(), a.Value, c.typID(n.Type())))
	return &Code{
		Pre:   pre,
		Value: fmt.Sprintf("(*((%s*)(%s->value)))", c.CType(n.Type()), a.Value),
	}, nil
}
