package ir

import (
	"fmt"

	"github.com/Nv7-Github/bsharp/parser"
)

func (b *Builder) Build(p *parser.Parser, fs FS) error {
	b.imported[p.Filename()] = empty{} // Mark as imported

	err := b.importPass(p, fs)
	if err != nil {
		return err
	}
	err = b.defPass(p)
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

	if len(b.Errors) > 0 {
		return fmt.Errorf("built with %d errors", len(b.Errors))
	}
	return nil
}
