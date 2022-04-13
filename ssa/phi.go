package ssa

import (
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

type Phi struct {
	Values   []ID
	Typ      types.Type
	Variable int // Annotation for putting back mem
}

func (p *Phi) Type() types.Type { return p.Typ }
func (p *Phi) String() string {
	out := &strings.Builder{}
	out.WriteString("φ(")
	for i, val := range p.Values {
		out.WriteString(val.String())
		if i != len(p.Values)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteRune(')')
	return out.String()
}
func (p *Phi) Args() []ID     { return p.Values }
func (p *Phi) SetArgs(v []ID) { p.Values = v }
