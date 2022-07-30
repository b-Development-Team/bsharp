package bot

import (
	"errors"
	"strings"
	"time"

	"github.com/Nv7-Github/bsharp/backends/cgen"
	"github.com/Nv7-Github/bsharp/backends/interpreter"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/sevcord"
)

type fs struct {
	b   *Bot
	gld string
}

func (f *fs) Parse(file string) (*parser.Parser, error) {
	id := strings.TrimSuffix(file, ".bsp")
	dat, err := f.b.Get(f.gld)
	if err != nil {
		return nil, err
	}
	src, rsp := dat.GetSource(id)
	if !rsp.Suc {
		return nil, errors.New(rsp.Msg)
	}
	stream := tokens.NewStream(file, src)
	tok := tokens.NewTokenizer(stream)
	err = tok.Tokenize()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(tok)
	err = parser.Parse()
	if err != nil {
		return nil, err
	}
	return parser, nil
}

// ctxWriter is an io.Writer implementation around *Ctx
type ctxWriter struct {
	sevcord.Ctx
	lastSent time.Time
	data     *strings.Builder
	prog     *interpreter.Interpreter
	cmps     [][]sevcord.Component
}

func newCtxWriter(ctx sevcord.Ctx) *ctxWriter {
	return &ctxWriter{Ctx: ctx}
}

func (c *ctxWriter) Setup() {
	emb := sevcord.NewEmbedBuilder("Program Output").Color(3447003).Description("Running...")
	rsp := sevcord.EmbedResponse(emb)
	for _, cmp := range c.cmps {
		rsp.ComponentRow(cmp...)
	}

	c.data = &strings.Builder{}
	c.Respond(rsp)
	c.lastSent = time.Now()
}

func (c *ctxWriter) Flush() error {
	var err error
	var emb *sevcord.EmbedBuilder
	if len(c.data.String()) == 0 {
		emb = sevcord.NewEmbedBuilder("Program Successful").Color(5763719).Description("⚠️ **WARNING**: No program output!")
	} else {
		emb = sevcord.NewEmbedBuilder("Program Output").Color(5763719).Description("```" + c.data.String() + "```")
	}
	rsp := sevcord.EmbedResponse(emb)
	for _, cmp := range c.cmps {
		rsp.ComponentRow(cmp...)
	}

	c.Edit(rsp)
	c.lastSent = time.Now()
	return err
}

func (c *ctxWriter) Write(b []byte) (int, error) {
	n, err := c.data.Write(b)
	if time.Since(c.lastSent).Milliseconds() > 200 { // Update every 200ms
		err = c.Flush()
	}
	return n, err
}

func (c *ctxWriter) Error(err error) {
	var emb *sevcord.EmbedBuilder
	if len(c.data.String()) == 0 {
		emb = sevcord.NewEmbedBuilder("Runtime Error").Color(15548997).Field("Error", "```"+err.Error()+"```", false)
	} else {
		emb = sevcord.NewEmbedBuilder("Runtime Error").Color(15548997).Field("Output", "```"+c.data.String()+"```", false).Field("Error", "```"+err.Error()+"```", false)
	}

	rsp := sevcord.EmbedResponse(emb)
	for _, cmp := range c.cmps {
		rsp.ComponentRow(cmp...)
	}

	c.Edit(rsp)
}

func (b *Bot) BuildCode(filename string, src string, ctx sevcord.Ctx) (*ir.IR, error) {
	stream := tokens.NewStream(filename, src)
	tok := tokens.NewTokenizer(stream)
	err := tok.Tokenize()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(tok)
	err = parser.Parse()
	if err != nil {
		return nil, err
	}
	bld := ir.NewBuilder()
	for _, ext := range Exts {
		bld.AddExtension(ext)
	}
	err = bld.Build(parser, &fs{b: b, gld: ctx.Guild()})
	if err != nil {
		e := &strings.Builder{}
		e.WriteString("```")
		for _, v := range bld.Errors {
			e.WriteString(v.Pos.Error("%s", v.Message).Error())
			e.WriteRune('\n')
		}
		e.WriteString(err.Error())
		e.WriteString("```")
		return nil, errors.New(e.String())
	}
	return bld.IR(), nil
}

func (b *Bot) CompileCode(filename, src string, ctx sevcord.Ctx) (string, error) {
	ir, err := b.BuildCode(filename, src, ctx)
	if err != nil {
		return "", err
	}
	cgen := cgen.NewCGen(ir)
	return cgen.Build()
}

const MaxTime = time.Second * 150

func (c *ctxWriter) UpdateCmps() {
	c.cmps = [][]sevcord.Component{{
		&sevcord.Button{
			Label: "Stop",
			Style: sevcord.ButtonStyleDanger,
			Handler: func(ctx sevcord.Ctx) {
				if ctx.User().ID != c.User().ID {
					return
				}
				c.prog.Stop("program stopped by user")
				c.Ctx = ctx
			},
		},
	}}
}

func (b *Bot) RunCode(filename string, src string, ctx sevcord.Ctx, extensionCtx *extensionCtx) error {
	ir, err := b.BuildCode(filename, src, ctx)
	if err != nil {
		return err
	}
	stdout := newCtxWriter(ctx)
	stdout.UpdateCmps()
	stdout.Setup()
	extensionCtx.Stdout = stdout
	interp := interpreter.NewInterpreter(ir, stdout)
	stdout.prog = interp

	exts := getExtensions(extensionCtx)
	for _, ext := range exts {
		interp.AddExtension(ext)
	}

	// Run
	var done = false
	go func() {
		time.Sleep(MaxTime)
		if !done {
			interp.Stop("program exceeded time limit of 2.5 minutes")
		}
	}()

	err = interp.Run()
	stdout.cmps = [][]sevcord.Component{{
		&sevcord.Button{
			Label:    "Stop",
			Style:    sevcord.ButtonStyleDanger,
			Disabled: true,
		},
	}}
	done = true
	flusherr := stdout.Flush()
	if err != nil {
		stdout.Error(err)
		return nil
	}
	return flusherr
}
