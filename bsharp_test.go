package main

import (
	"embed"
	"errors"
	"io"
	"path/filepath"
	"testing"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
)

//go:embed std/*.bsp examples/*.bsp
var fuzzcorpus embed.FS

type FuzzFS struct{}

func (f *FuzzFS) Parse(name string) (*parser.Parser, error) {
	return nil, errors.New("not implemented")
}

func FuzzIR(f *testing.F) {
	files, err := fuzzcorpus.ReadDir("std")
	if err != nil {
		f.Fatal(err)
	}
	for _, fi := range files {
		v, err := fuzzcorpus.Open(filepath.Join("std", fi.Name()))
		if err != nil {
			f.Fatal(err)
		}
		c, err := io.ReadAll(v)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(c))
		v.Close()
	}
	files, err = fuzzcorpus.ReadDir("examples")
	if err != nil {
		f.Fatal(err)
	}
	for _, fi := range files {
		v, err := fuzzcorpus.Open(filepath.Join("examples", fi.Name()))
		if err != nil {
			f.Fatal(err)
		}
		c, err := io.ReadAll(v)
		if err != nil {
			f.Fatal(err)
		}
		f.Add(string(c))
		v.Close()
	}

	f.Fuzz(func(t *testing.T, code string) {
		s := tokens.NewStream("main.bsp", code)
		tok := tokens.NewTokenizer(s)
		err := tok.Tokenize()
		if err != nil {
			return
		}

		p := parser.NewParser(tok)
		err = p.Parse()
		if err != nil {
			return
		}

		bld := ir.NewBuilder()
		err = bld.Build(p, &FuzzFS{})
		if err != nil {
			return
		}
	})
}
