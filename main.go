package main

import (
	_ "embed"
	"fmt"

	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
)

//go:embed examples/main.bsp
var code string

func main() {
	stream := tokens.NewStream("main.bsp", code)
	tok := tokens.NewTokenizer(stream)
	tok.Tokenize()

	parser := parser.NewParser(tok)
	err := parser.Parse()
	if err != nil {
		panic(err)
	}
	for _, node := range parser.Nodes {
		fmt.Println(node)
	}
}
