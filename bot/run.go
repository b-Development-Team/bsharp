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
	"github.com/bwmarrin/discordgo"
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
	*Ctx
	lastSent time.Time
	data     *strings.Builder
	prog     *interpreter.Interpreter
	cmps     []discordgo.MessageComponent
}

func newCtxWriter(ctx *Ctx) *ctxWriter {
	return &ctxWriter{Ctx: ctx}
}

var runCmp = []discordgo.MessageComponent{
	discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Stop",
				Style:    discordgo.DangerButton,
				CustomID: "stop",
			},
		},
	},
}

var runCmpEnd = []discordgo.MessageComponent{
	discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Stop",
				Style:    discordgo.DangerButton,
				CustomID: "stop",
				Disabled: true,
			},
		},
	},
}

func (c *ctxWriter) AddBtnHandler() {
	c.BtnHandler(func(data discordgo.MessageComponentInteractionData, ctx *Ctx) {
		if ctx.i.Member.User.ID != c.i.Member.User.ID {
			return
		}
		c.prog.Stop("program stopped by user")
		c.i = ctx.i
		c.isButton = true
	})
}

func (c *ctxWriter) Setup() {
	c.Followup()
	c.data = &strings.Builder{}
	c.Embed(&discordgo.MessageEmbed{
		Title:       "Program Output",
		Color:       3447003, // Blue
		Description: "Running...",
	}, c.cmps...)
	c.AddBtnHandler()
	c.lastSent = time.Now()
}

func (c *ctxWriter) Flush() error {
	var err error
	if len(c.data.String()) == 0 {
		err = c.Embed(&discordgo.MessageEmbed{
			Title:       "Program Successful",
			Color:       5763719, // Green
			Description: "⚠️ **WARNING**: No program output!",
		}, c.cmps...)
	} else {
		err = c.Embed(&discordgo.MessageEmbed{
			Title:       "Program Output",
			Color:       5763719, // Green
			Description: "```" + c.data.String() + "```",
		}, c.cmps...)
	}
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

func (c *ctxWriter) Error(err error) error {
	if len(c.data.String()) == 0 {
		return c.Embed(&discordgo.MessageEmbed{
			Title: "Runtime Error",
			Color: 15548997, // Red
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Error",
					Value: "```" + err.Error() + "```",
				},
			},
		}, c.cmps...)
	}
	return c.Embed(&discordgo.MessageEmbed{
		Title: "Runtime Error",
		Color: 15548997, // Red
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Output",
				Value: "```" + c.data.String() + "```",
			},
			{
				Name:  "Error",
				Value: "```" + err.Error() + "```",
			},
		},
	}, c.cmps...)
}

func (b *Bot) BuildCode(filename string, src string, ctx *Ctx) (*ir.IR, error) {
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
		for _, v := range bld.Errors {
			e.WriteString(v.Pos.Error("%s", v.Message).Error() + "\n")
		}
		e.WriteString(err.Error())
		return nil, errors.New(e.String())
	}
	return bld.IR(), nil
}

func (b *Bot) CompileCode(filename, src string, ctx *Ctx) (string, error) {
	ir, err := b.BuildCode(filename, src, ctx)
	if err != nil {
		return "", err
	}
	cgen := cgen.NewCGen(ir)
	return cgen.Build()
}

const MaxTime = time.Second * 150

func (b *Bot) RunCode(filename string, src string, ctx *Ctx, extensionCtx *extensionCtx) error {
	ir, err := b.BuildCode(filename, src, ctx)
	if err != nil {
		return err
	}
	stdout := newCtxWriter(ctx)
	stdout.cmps = runCmp
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
	stdout.cmps = runCmpEnd
	done = true
	flusherr := stdout.Flush()
	if err != nil {
		err = stdout.Error(err)
		if err != nil {
			return err
		}
	}
	return flusherr
}
