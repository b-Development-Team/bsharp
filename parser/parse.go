package parser

import (
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
		if t.Value == "TRUE" || t.Value == "FALSE" {
			return &BoolNode{
				pos:   t.Pos,
				Value: t.Value == "TRUE",
			}, nil
		}
		if t.Value == "NULL" {
			return &NullNode{
				pos: t.Pos,
			}, nil
		}
		return &IdentNode{
			pos:   t.Pos,
			Value: t.Value,
		}, nil

	case tokens.TokenTypeNumber:
		t := p.t.Tok()
		p.t.Eat()
		return &NumberNode{
			pos:   t.Pos,
			Value: t.Value,
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
		if p.t.Tok().Typ == tokens.TokenTypeRBrack { // Eat ]
			pos = pos.Extend(p.t.Tok().Pos)
			p.t.Eat()
			break
		}

		arg, err := p.ParseNode()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	return &CallNode{
		pos:  pos,
		Name: name,
		Args: args,
	}, nil
}
