package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/Nv7-Github/bsharp/backends/cgen"
	"github.com/Nv7-Github/bsharp/backends/interpreter"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/alexflint/go-arg"
)

type Run struct {
	Files []string `arg:"positional,-i,--input" help:"input B# program"`
}

type Build struct {
	Files []string `arg:"positional,-i,--input" help:"input B# program"`
}

type Args struct {
	Run   *Run `arg:"subcommand:run" help:"run a B# program"`
	Build *Run `arg:"subcommand:build" help:"compile a B# program"`
	Time  bool `help:"print timing for each stage" arg:"-t"`
}

func main() {
	args := Args{}
	p := arg.MustParse(&args)

	switch {
	case args.Run != nil:
		start := time.Now()
		files := make(map[string]struct{}, len(args.Run.Files))
		for _, f := range args.Run.Files {
			files[f] = struct{}{}
		}

		// Build
		fs := &dirFS{files}
		v, err := fs.Parse(args.Run.Files[0])
		if err != nil {
			p.Fail(err.Error())
		}
		ir := ir.NewBuilder()
		err = ir.Build(v, fs)
		if err != nil {
			p.Fail(err.Error())
		}

		if args.Time {
			fmt.Println("Built in", time.Since(start))
		}

		// Run
		start = time.Now()
		interp := interpreter.NewInterpreter(ir.IR(), os.Stdout)
		err = interp.Run()
		if err != nil {
			p.Fail(err.Error())
		}
		if args.Time {
			fmt.Println("Ran in", time.Since(start))
		}

	case args.Build != nil:
		start := time.Now()
		files := make(map[string]struct{}, len(args.Build.Files))
		for _, f := range args.Build.Files {
			files[f] = struct{}{}
		}

		// Build
		fs := &dirFS{files}
		v, err := fs.Parse(args.Build.Files[0])
		if err != nil {
			p.Fail(err.Error())
		}
		ir := ir.NewBuilder()
		err = ir.Build(v, fs)
		if err != nil {
			p.Fail(err.Error())
		}

		if args.Time {
			fmt.Println("Built in", time.Since(start))
		}

		// Run
		start = time.Now()
		cgen := cgen.NewCGen(ir.IR())
		code, err := cgen.Build()
		if err != nil {
			p.Fail(err.Error())
		}
		if args.Time {
			fmt.Println("Ran in", time.Since(start))
		}

		// TODO: Compile and save, put C code in tmp file
		fmt.Println(code)
	}
}
