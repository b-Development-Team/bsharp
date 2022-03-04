package ir

import (
	"strconv"
	"strings"

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

func NewConst(typ types.Type, pos *tokens.Pos, val interface{}) *Const {
	return &Const{
		typ:   typ,
		pos:   pos,
		Value: val,
	}
}

func (b *Builder) buildNumber(n *parser.NumberNode) (Node, error) {
	if strings.Contains(n.Value, ".") {
		v, err := strconv.ParseFloat(n.Value, 64)
		if err != nil {
			return nil, err
		}
		return &Const{
			typ:   types.FLOAT,
			pos:   n.Pos(),
			Value: v,
		}, nil
	}
	v, err := strconv.ParseInt(n.Value, 10, 64)
	if err != nil {
		return nil, err
	}
	return &Const{
		typ:   types.INT,
		pos:   n.Pos(),
		Value: int(v),
	}, nil
}
