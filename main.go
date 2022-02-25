package main

import (
	_ "embed"
	"fmt"

	"github.com/Nv7-Github/bsharp/tokens"
)

//go:embed examples/main.bsharp
var code string

func main() {
	stream := tokens.NewStream("main.bsharp", code)
	tok := tokens.NewTokenizer(stream)
	tok.Tokenize()
	for _, tok := range tok.Tokens {
		fmt.Println(tok)
	}
}
