package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
)

type Const struct {
	typ types.Type
	pos *tokens.Pos

	Value interface{}
}

func (c *Const) Type() types.Type { return c.typ }
func (c *Const) Pos() *tokens.Pos { return c.pos }

func (b *Builder) buildString(n *parser.StringNode) Node {
	return &Const{
		typ:   types.STRING,
		pos:   n.Pos(),
		Value: n.Value,
	}
}

func (b *Builder) buildIdent(n *parser.IdentNode) Node {
	return &Const{
		typ:   types.IDENT,
		pos:   n.Pos(),
		Value: n.Value,
	}
}

func (b *Builder) buildNumber(n *parser.NumberNode) Node {
	if float64(int(n.Value)) == n.Value {
		return &Const{
			typ:   types.INT,
			pos:   n.Pos(),
			Value: int(n.Value),
		}
	}
	return &Const{
		typ:   types.FLOAT,
		pos:   n.Pos(),
		Value: n.Value,
	}
}
