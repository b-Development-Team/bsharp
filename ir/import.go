package ir

import "github.com/Nv7-Github/bsharp/parser"

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
			return call.Pos().Error("file \"%s\" already imported!", name)
		}

		// Get file
		file, err := fs.Parse(name)
		if err != nil {
			return call.Pos().Error("%s", err.Error())
		}

		// Build file
		err = b.Build(file, fs)
		if err != nil {
			return err
		}
	}

	return nil
}
