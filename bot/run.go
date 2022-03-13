package bot

import (
	"errors"
	"strings"
	"time"

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
}

func newCtxWriter(ctx *Ctx) *ctxWriter {
	return &ctxWriter{Ctx: ctx}
}

func (c *ctxWriter) Setup() {
	c.Followup()
	c.data = &strings.Builder{}
	c.Embed(&discordgo.MessageEmbed{
		Title:       "Program Output",
		Color:       3447003, // Blue
		Description: "Running...",
	})
	c.lastSent = time.Now()
}

func (c *ctxWriter) Flush() error {
	var err error
	if len(c.data.String()) == 0 {
		err = c.Embed(&discordgo.MessageEmbed{
			Title:       "Program Successful",
			Color:       5763719, // Green
			Description: "⚠️ **WARNING**: No program output!",
		})
	} else {
		err = c.Embed(&discordgo.MessageEmbed{
			Title:       "Program Output",
			Color:       5763719, // Green
			Description: "```" + c.data.String() + "```",
		})
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
		})
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
	})
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
	for _, ext := range exts {
		bld.AddExtension(ext)
	}
	err = bld.Build(parser, &fs{b: b, gld: ctx.Guild()})
	if err != nil {
		return nil, err
	}
	return bld.IR(), nil
}

func (b *Bot) RunCode(filename string, src string, ctx *Ctx, extensionCtx *extensionCtx) error {
	ir, err := b.BuildCode(filename, src, ctx)
	if err != nil {
		return err
	}
	stdout := newCtxWriter(ctx)
	stdout.Setup()
	interp := interpreter.NewInterpreter(ir, stdout)

	exts := getExtensions(extensionCtx)
	for _, ext := range exts {
		interp.AddExtension(ext)
	}

	// Run
	err = interp.Run()
	flusherr := stdout.Flush()
	if err != nil {
		err = stdout.Error(err)
		if err != nil {
			return err
		}
	}
	return flusherr
}
