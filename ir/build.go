package ir

import "github.com/Nv7-Github/bsharp/parser"

func (b *Builder) Build(p *parser.Parser, fs FS) error {
	b.imported[p.Filename()] = empty{} // Mark as imported

	err := b.importPass(p, fs)
	if err != nil {
		return err
	}
	err = b.functionPass(p)
	if err != nil {
		return err
	}

	// Build nodes
	for _, node := range p.Nodes {
		node, err := b.buildNode(node)
		if err != nil {
			return err
		}
		if node != nil {
			b.Body = append(b.Body, node)
		}
	}

	return nil
}
