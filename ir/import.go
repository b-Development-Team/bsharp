package ir

import (
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/std"
	"github.com/Nv7-Github/bsharp/tokens"
)

type FS interface {
	Parse(src string) (*parser.Parser, error)
}

func (b *Builder) importPass(p *parser.Parser, fs FS) error {
	for _, node := range p.Nodes {
		call, ok := node.(*parser.CallNode)
		if !ok {
			continue
		}
		if call.Name != "IMPORT" {
			continue
		}

		// Its an import!
		if len(call.Args) != 1 {
			return call.Pos().Error("expected 1 argument to IMPORT")
		}

		// Get name
		nm := call.Args[0]
		nameV, ok := nm.(*parser.StringNode)
		if !ok {
			return nm.Pos().Error("expected import name")
		}
		name := nameV.Value

		// Check if imported
		_, exists := b.imported[name]
		if exists {
			continue
		}

		// Get file
		var p *parser.Parser
		_, exists = std.Std[name]
		if exists {
			// Parse
			stream := tokens.NewStream(name, std.Std[name])
			tok := tokens.NewTokenizer(stream)
			err := tok.Tokenize()
			if err != nil {
				return call.Pos().Error("%s", err.Error())
			}
			p = parser.NewParser(tok)
			err = p.Parse()
			if err != nil {
				return err
			}
		} else {
			var err error
			p, err = fs.Parse(name)
			if err != nil {
				return call.Pos().Error("%s", err.Error())
			}
		}

		// Build file
		err := b.Build(p, fs)
		if err != nil {
			return err
		}
	}

	return nil
}
