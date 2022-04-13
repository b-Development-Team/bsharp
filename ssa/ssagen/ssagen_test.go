package ssagen

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/ssa/constrm"
	"github.com/Nv7-Github/bsharp/ssa/dce"
	"github.com/Nv7-Github/bsharp/ssa/memrm"
	"github.com/Nv7-Github/bsharp/tokens"
)

const code = `# SSAGen Test
[DEFINE i 0]

#[IF [COMPARE [VAR i] == 0]
#  [DEFINE i 1]
#]

#[IF [COMPARE [VAR i] == 1]
#  [DEFINE i 2]
#ELSE
#  [DEFINE i 2]
#]

[WHILE [COMPARE [VAR i] < 10]
  [DEFINE i [MATH [VAR i] + 1]]
  [PRINT [STRING [VAR i]]]
]

[PRINT [STRING [VAR i]]]
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

	// Memrm Pass
	memrm := memrm.NewMemRM(s)
	memrm.Eval()

	fmt.Println("BEFORE OPTIMIZATION:")
	fmt.Println(s)

	// Constant folding
	constrm.Constrm(s)
	constrm.Phirm(s)

	// Dead code elimination
	dce := dce.NewDCE(s)
	dce.Remove()

	fmt.Println("AFTER:")
	fmt.Println(s)
}
