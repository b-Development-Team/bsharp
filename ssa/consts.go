package ssa

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
)

type Concat struct {
	Values []ID
}

func (c *Concat) Type() types.Type { return types.STRING }
func (c *Concat) String() string {
	out := &strings.Builder{}
	out.WriteString("Concat (")
	for i, val := range c.Values {
		out.WriteString(val.String())
		if i != len(c.Values)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteRune(')')
	return out.String()
}
func (c *Concat) Args() []ID     { return c.Values }
func (c *Concat) SetArgs(v []ID) { c.Values = v }

type LogicalOp struct {
	Op  ir.LogicalOp
	Lhs ID
	Rhs *ID // nil if no rhs
}

func (l *LogicalOp) Type() types.Type { return types.BOOL }
func (l *LogicalOp) String() string {
	if l.Rhs != nil {
		return fmt.Sprintf("%s (%s, %s)", l.Op.String(), l.Lhs.String(), l.Rhs.String())
	}
	return fmt.Sprintf("%s %s", l.Op.String(), l.Lhs.String())
}
func (l *LogicalOp) Args() []ID {
	if l.Rhs != nil {
		return []ID{l.Lhs, *l.Rhs}
	}
	return []ID{l.Lhs}
}
func (l *LogicalOp) SetArgs(v []ID) {
	l.Lhs = v[0]
	if len(v) > 1 {
		l.Rhs = &v[1]
	}
}
