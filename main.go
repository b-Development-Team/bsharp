package main

import (
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Nv7-Github/bsharp/backends/interpreter"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
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

func extadd(a, b int) int {
	return a + b
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

	// Extensions
	ext, err := interpreter.NewExtension(extadd, "EXTADD")
	if err != nil {
		panic(err)
	}

	builder := ir.NewBuilder()
	builder.AddExtension(ext.IRExtension())
	err = builder.Build(parser, &dirFS{dir: "examples"})
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	interp := interpreter.NewInterpreter(builder.IR(), os.Stdout)
	interp.AddExtension(ext)
	fmt.Println(ext.Call([]interface{}{1, 2}))
	err = interp.Run()
	if err != nil {
		panic(err)
	}
}
