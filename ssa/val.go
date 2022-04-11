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

type Cast struct {
	Value ID
	From  types.Type
	To    types.Type
}

func (c *Cast) Type() types.Type { return c.To }
func (c *Cast) String() string {
	return fmt.Sprintf("Cast (%s)[%s -> %s]", c.Value.String(), c.From.String(), c.To.String())
}
