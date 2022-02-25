package parser

import "github.com/Nv7-Github/bsharp/tokens"

type Parser struct {
	t     *tokens.Tokenizer
	Nodes []Node
}

func NewParser(t *tokens.Tokenizer) *Parser {
	return &Parser{
		t:     t,
		Nodes: make([]Node, 0),
	}
}

func (p *Parser) Parse() error {
	for p.t.HasNext() {
		n, err := p.ParseNode()
		if err != nil {
			return err
		}
		p.Nodes = append(p.Nodes, n)
	}
	return nil
}
