package ssagen

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/ssa/phirm"
	"github.com/Nv7-Github/bsharp/ssa/pipeline"
	"github.com/Nv7-Github/bsharp/tokens"
)

const code = `# SSAGen Test
[DEFINE a ""]
[SWITCH [VAR a]
	[CASE "Hi"
		[DEFINE a "e"]
	]

	[CASE "Hey"
		[DEFINE a "b"]
	]

	[DEFAULT
		[DEFINE a "h"]
	]
]
[PRINT [VAR a]]
[PRINT "Hi"]
`

type fs struct{}

func (*fs) Parse(f string) (*parser.Parser, error) {
	return nil, errors.New("not implemented")
}

func TestSSAGen(t *testing.T) {
	stream := tokens.NewStream("main.bsp", code)
	tok := tokens.NewTokenizer(stream)
	err := tok.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	parse := parser.NewParser(tok)
	err = parse.Parse()
	if err != nil {
		t.Fatal(err)
	}
	ir := ir.NewBuilder()
	err = ir.Build(parse, &fs{})
	if err != nil {
		t.Fatal(err)
	}

	// Actually generate
	i := ir.IR()
	gen := NewSSAGen(i)
	gen.Build()
	s := gen.SSA()

	fmt.Println("BEFORE:")
	fmt.Println(s)

	p := pipeline.New()
	p.ConstantPropagation()
	p.DeadCodeElimination()
	p.Run(s)

	fmt.Println("AFTER:")
	fmt.Println(s)

	rm := phirm.NewPhiRM(s)
	rm.Remove()
	fmt.Println("PHIRM:")
	fmt.Println(s)

	// Rebuild B#
	/*g := bspgen.NewBSPGen(s, i)
	out := g.Build()
	fmt.Println("B#:")
	spew.Dump(out.Body)*/
}
