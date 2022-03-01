package main

import (
	"fmt"
	"os"

	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
)

type dirFS struct {
	files map[string]struct{}
}

func (d *dirFS) Parse(name string) (*parser.Parser, error) {
	if _, ok := d.files[name]; !ok {
		return nil, fmt.Errorf("bsharp: file not found: %s", name)
	}
	src, err := os.ReadFile(name)
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
