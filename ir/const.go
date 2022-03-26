package ir

import (
	"fmt"
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
func (c *Const) Code(cnf CodeConfig) string {
	switch c.typ.BasicType() {
	case types.INT:
		return fmt.Sprintf("%d", c.Value.(int))

	case types.FLOAT:
		return strconv.FormatFloat(c.Value.(float64), 'f', -1, 64)

	case types.STRING:
		return fmt.Sprintf("%q", c.Value.(string))

	case types.IDENT:
		return c.Value.(string)

	case types.BOOL:
		if c.Value.(bool) {
			return "TRUE"
		}
		return "FALSE"

	default:
		return "[UNKNOWN CONST]"
	}
}

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

func (b *Builder) buildBool(n *parser.BoolNode) Node {
	return &Const{
		typ:   types.BOOL,
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
