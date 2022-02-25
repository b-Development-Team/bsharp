package parser

import (
	"strconv"

	"github.com/Nv7-Github/bsharp/tokens"
)

func (p *Parser) ParseNode() (Node, error) {
	// Check if not call node
	switch p.t.Tok().Typ {
	case tokens.TokenTypeString:
		t := p.t.Tok()
		p.t.Eat()
		return &StringNode{
			pos:   t.Pos,
			Value: t.Value,
		}, nil

	case tokens.TokenTypeIdent:
		t := p.t.Tok()
		p.t.Eat()
		return &IdentNode{
			pos:   t.Pos,
			Value: t.Value,
		}, nil

	case tokens.TokenTypeNumber:
		v, err := strconv.ParseFloat(p.t.Tok().Value, 64)
		if err != nil {
			return nil, p.t.Tok().Pos.Error("invalid number")
		}
		t := p.t.Tok()
		p.t.Eat()
		return &NumberNode{
			pos:   t.Pos,
			Value: v,
		}, nil
	}

	// LBrack
	pos := p.t.Tok().Pos
	if p.t.Tok().Typ != tokens.TokenTypeLBrack {
		return nil, p.t.Tok().Pos.Error("expected '['")
	}
	p.t.Eat()

	// Name
	if !p.t.HasNext() || p.t.Tok().Typ != tokens.TokenTypeIdent {
		return nil, p.t.Last().Error("expected identifier")
	}
	name := p.t.Tok().Value
	p.t.Eat()

	// Args
	args := make([]Node, 0)
	for p.t.HasNext() {
		arg, err := p.ParseNode()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)

		if p.t.Tok().Typ == tokens.TokenTypeRBrack { // Eat ]
			p.t.Eat()
			break
		}
	}

	return &CallNode{
		pos:  pos,
		Name: name,
		Args: args,
	}, nil
}
