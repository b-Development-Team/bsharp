package optimize

import (
	"errors"
	"testing"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/davecgh/go-spew/spew"
)

const code = `[DEFINE a [STRING [MATH 1 + 1]]]
[PRINT [VAR a]]
`

type fs struct{}

func (*fs) Parse(f string) (*parser.Parser, error) {
	return nil, errors.New("not implemented")
}

func TestOptimizer(t *testing.T) {
	stream := tokens.NewStream("main.bsp", code)
	tok := tokens.NewTokenizer(stream)
	tok.Tokenize()
	parse := parser.NewParser(tok)
	parse.Parse()
	ir := ir.NewBuilder()
	ir.Build(parse, &fs{})

	opt := NewOptimizer(ir.IR())
	out := opt.Optimize()
	spew.Dump(out)
}
