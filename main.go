package main

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/davecgh/go-spew/spew"
)

//go:embed examples/main.bsp
var code string

type dirFS struct {
	dir string
}

func (d *dirFS) Parse(name string) (*parser.Parser, error) {
	src, err := os.ReadFile(filepath.Join(d.dir, name))
	if err != nil {
		return nil, err
	}
	stream := tokens.NewTokenizer(tokens.NewStream(name, string(src)))
	err = stream.Tokenize()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(stream)

	err = parser.Parse()
	return parser, err
}

func main() {
	stream := tokens.NewStream("main.bsp", code)
	tok := tokens.NewTokenizer(stream)
	err := tok.Tokenize()
	if err != nil {
		panic(err)
	}

	parser := parser.NewParser(tok)
	err = parser.Parse()
	if err != nil {
		panic(err)
	}

	builder := ir.NewBuilder()
	err = builder.Build(parser, &dirFS{dir: "examples"})
	if err != nil {
		panic(err)
	}

	for _, fn := range builder.Funcs {
		spew.Dump(fn)
		fmt.Println()
		fmt.Println()
	}

	for _, node := range builder.Body {
		spew.Dump(node)
	}
}
