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
	return nil
}
