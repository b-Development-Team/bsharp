package ssa

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Const struct {
	Value interface{}
	Typ   types.Type
}

func (c *Const) Type() types.Type { return c.Typ }
func (c *Const) String() string {
	return fmt.Sprintf("Const {%s}(%v)", c.Type().String(), c.Value)
}
func (c *Const) Args() []ID     { return []ID{} }
func (c *Const) SetArgs(_ []ID) {}

type Compare struct {
	Op  ir.CompareOperation
	Lhs ID
	Rhs ID
	Typ types.Type
}

func (c *Compare) Type() types.Type { return c.Typ }
func (c *Compare) String() string {
	return fmt.Sprintf("Compare (%s) %s (%s)", c.Lhs.String(), c.Op.String(), c.Rhs.String())
}
func (c *Compare) Args() []ID { return []ID{c.Lhs, c.Rhs} }
func (c *Compare) SetArgs(v []ID) {
	c.Lhs = v[0]
	c.Rhs = v[1]
}

type Math struct {
	Op  ir.MathOperation
	Lhs ID
	Rhs ID
	Typ types.Type
}

func (m *Math) Type() types.Type { return m.Typ }
func (m *Math) String() string {
	return fmt.Sprintf("Math (%s) %s (%s)", m.Lhs.String(), m.Op.String(), m.Rhs.String())
}
func (m *Math) Args() []ID { return []ID{m.Lhs, m.Rhs} }
func (m *Math) SetArgs(v []ID) {
	m.Lhs = v[0]
	m.Rhs = v[1]
}

type Cast struct {
	Value ID
	From  types.Type
	To    types.Type
}

func (c *Cast) Type() types.Type { return c.To }
func (c *Cast) String() string {
	return fmt.Sprintf("Cast (%s)[%s -> %s]", c.Value.String(), c.From.String(), c.To.String())
}
func (c *Cast) Args() []ID { return []ID{c.Value} }
func (c *Cast) SetArgs(v []ID) {
	c.Value = v[0]
}
