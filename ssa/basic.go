package ssa

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/types"
)

type Print struct {
	Value ID
}

func (p *Print) Type() types.Type { return types.INT }
func (p *Print) String() string   { return fmt.Sprintf("Print (%s)", p.Value.String()) }
func (p *Print) Args() []ID       { return []ID{p.Value} }
func (p *Print) SetArgs(v []ID)   { p.Value = v[0] }

type Phi struct {
	Values []ID
	Typ    types.Type
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
