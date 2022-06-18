package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Nv7-Github/bsharp/backends/bsp"
	"github.com/Nv7-Github/bsharp/backends/bstar"
	"github.com/Nv7-Github/bsharp/backends/cgen"
	"github.com/Nv7-Github/bsharp/backends/interpreter"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/types"
	"github.com/alexflint/go-arg"
)

type Run struct {
	Files []string `arg:"positional,-i,--input" help:"input B# program"`
}

type Build struct {
	Files    []string `arg:"positional,-i,--input" help:"input B# program"`
	Output   string   `arg:"required,-o,--output" help:"output executable"`
	Optimize bool     `arg:"-O,--optimize" help:"whether to optimize during compiling"`
}

type BSPGen struct {
	Files  []string `arg:"positional,-i,--input" help:"input B# program"`
	Output string   `arg:"required,-o,--output" help:"output B# file"`
}

type BStarGen struct {
	Files  []string `arg:"positional,-i,--input" help:"input B# program"`
	Output string   `arg:"required,-o,--output" help:"output B# file"`
}

type Args struct {
	Run      *Run      `arg:"subcommand:run" help:"run a B# program"`
	Build    *Build    `arg:"subcommand:build" help:"compile a B# program"`
	BSPGen   *BSPGen   `arg:"subcommand:ir" help:"view the IR in B# form"`
	BStarGen *BStarGen `arg:"subcommand:bstar" help:"compile the B# program to B*"`

	Time bool `help:"print timing for each stage" arg:"-t"`
}

var exts = []*ir.Extension{
	{
		Name:    "INPUT",
		Params:  []types.Type{types.STRING},
		RetType: types.STRING,
	},
}

func build(args []string, ti bool, p *arg.Parser) *ir.IR {
	start := time.Now()
	files := make(map[string]struct{}, len(args))
	for _, f := range args {
		files[f] = struct{}{}
	}

	// Build
	fs := &dirFS{files}
	v, err := fs.Parse(args[0])
	if err != nil {
		p.Fail(err.Error())
	}
	ir := ir.NewBuilder()
	for _, ext := range exts {
		ir.AddExtension(ext)
	}
	err = ir.Build(v, fs)
	if err != nil {
		if len(ir.Errors) > 0 {
			for _, err := range ir.Errors {
				fmt.Println(err.Pos.Error("%s", err.Message))
			}
		}
		p.Fail(err.Error())
	}

	if ti {
		fmt.Println("Built in", time.Since(start))
	}

	return ir.IR()
}

func main() {
	args := Args{}
	p := arg.MustParse(&args)

	switch {
	case args.Run != nil:
		ir := build(args.Run.Files, args.Time, p)

		// Run
		start := time.Now()
		interp := interpreter.NewInterpreter(ir, os.Stdout)
		reader := bufio.NewReader(os.Stdin)
		interp.AddExtension(interpreter.NewExtension("INPUT", func(v []any) (any, error) {
			fmt.Print(v[0].(string))
			line, _, err := reader.ReadLine()
			if err != nil {
				return nil, err
			}
			return string(line), nil
		}, []types.Type{types.STRING}, types.STRING))
		// Cleanup function
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			interp.Stop("interrupted")
		}()

		err := interp.Run()
		if err != nil {
			p.Fail(err.Error())
		}
		if args.Time {
			fmt.Println("Ran in", time.Since(start))
		}

	case args.Build != nil:
		ir := build(args.Build.Files, args.Time, p)

		// CGen
		start := time.Now()
		cgen := cgen.NewCGen(ir)
		code, err := cgen.Build()
		if err != nil {
			p.Fail(err.Error())
		}
		if args.Time {
			fmt.Println("Generated code in", time.Since(start))
		}

		// Compile
		if strings.HasSuffix(args.Build.Output, ".c") {
			err = os.WriteFile(args.Build.Output, []byte(code), os.ModePerm)
			if err != nil {
				p.Fail(err.Error())
			}
			return
		}

		// Save code
		start = time.Now()
		f, err := os.CreateTemp("", "*.c")
		if err != nil {
			p.Fail(err.Error())
		}
		defer f.Close()
		_, err = f.WriteString(code)
		if err != nil {
			p.Fail(err.Error())
		}

		// Build
		o := "-O0"
		if args.Build.Optimize {
			o = "-O2"
		}
		cmd := exec.Command("cc", f.Name(), "-o", args.Build.Output, o)
		err = cmd.Run()
		if err != nil {
			p.Fail(err.Error())
		}

		if args.Time {
			fmt.Println("Compiled in", time.Since(start))
		}

	case args.BSPGen != nil:
		ir := build(args.BSPGen.Files, args.Time, p)

		// Make B#
		start := time.Now()
		gen := bsp.NewBSP(ir)
		out, err := gen.Build()
		if err != nil {
			p.Fail(err.Error())
		}
		if args.Time {
			fmt.Println("Generated B# in", time.Since(start))
		}

		// Save
		err = os.WriteFile(args.BSPGen.Output, []byte(out), os.ModePerm)
		if err != nil {
			p.Fail(err.Error())
		}

	case args.BStarGen != nil:
		exts[0] = &ir.Extension{
			Name:    "ARGS",
			Params:  []types.Type{types.INT},
			RetType: types.STRING,
		}
		ir := build(args.BStarGen.Files, args.Time, p)

		// Make B#
		start := time.Now()
		gen := bstar.NewBStar(ir)
		out, err := gen.Build()
		if err != nil {
			p.Fail(err.Error())
		}
		if args.Time {
			fmt.Println("Generated B# in", time.Since(start))
		}

		// Save
		conf := &bstar.BStarConfig{Seperator: " "}
		code := out.Code(conf)
		err = os.WriteFile(args.BStarGen.Output, []byte(code), os.ModePerm)
		if err != nil {
			p.Fail(err.Error())
		}
	}
}
