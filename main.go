package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
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
	Files    []string `arg:"positional,-i,--input" help:"input B# program"`
	Output   string   `arg:"required,-o,--output" help:"output executable"`
	Optimize bool     `arg:"-O,--optimize" help:"whether to optimize during compiling"`
}

type Args struct {
	Run   *Run   `arg:"subcommand:run" help:"run a B# program"`
	Build *Build `arg:"subcommand:build" help:"compile a B# program"`
	Time  bool   `help:"print timing for each stage" arg:"-t"`
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
			if len(ir.Errors) > 0 {
				for _, err := range ir.Errors {
					fmt.Println(err.Pos.Error("%s", err.Message))
				}
			}
			p.Fail(err.Error())
		}

		if args.Time {
			fmt.Println("Built in", time.Since(start))
		}

		// Run
		start = time.Now()
		interp := interpreter.NewInterpreter(ir.IR(), os.Stdout)
		// Cleanup function
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			interp.Stop("interrupted")
		}()

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
			if len(ir.Errors) > 0 {
				for _, err := range ir.Errors {
					fmt.Println(err.Pos.Error("%s", err.Message))
				}
			}
			p.Fail(err.Error())
		}

		if args.Time {
			fmt.Println("Built in", time.Since(start))
		}

		// CGen
		start = time.Now()
		cgen := cgen.NewCGen(ir.IR())
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
	}
}
