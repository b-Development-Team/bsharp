package optimize

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
)

const code = `[DEFINE a [STRING [MATH 1 + 1]]]
[PRINT [VAR a]]

[DEFINE i 0]
[WHILE [COMPARE [VAR i] < 10]
	[PRINT [STRING [VAR i]]]
	[DEFINE i [MATH [VAR i] + 1]]
]
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
	cnf := ir.CodeConfig{
		Indent: 4,
	}
	ir := ir.NewBuilder()
	ir.Build(parse, &fs{})

	i := ir.IR()
	fmt.Println("BEFORE: ")
	fmt.Println(i.Code(cnf))

	opt := NewOptimizer(i)
	out := opt.Optimize()

	fmt.Println("AFTER: ")
	fmt.Println(out.Code(cnf))
}
