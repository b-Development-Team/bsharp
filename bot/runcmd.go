package bot

import (
	"io"
	"net/http"
	"time"

	"github.com/Nv7-Github/sevcord"
)

func RunCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommandGroup{
		Name:        "run",
		Description: "Run a B# program!",
		Children: []sevcord.SlashCommandObject{
			&sevcord.SlashCommand{
				Name:        "code",
				Description: "Run the source code of a B# program!",
				Options:     []sevcord.Option{},
				Handler:     b.RunCodeCmd,
			},
			&sevcord.SlashCommand{
				Name:        "tag",
				Description: "Run a tag!",
				Options: []sevcord.Option{
					{
						Kind:         sevcord.OptionKindString,
						Name:         "tag",
						Description:  "The ID of the tag to run!",
						Required:     true,
						Autocomplete: b.Autocomplete,
					},
				},
				Handler: b.RunTagCmd,
			},
			&sevcord.SlashCommand{
				Name:        "file",
				Description: "Run the source code of a B# program, uploaded as a file!",
				Options: []sevcord.Option{
					{
						Kind:        sevcord.OptionKindAttachment,
						Name:        "file",
						Description: "The file to run!",
						Required:    true,
					},
				},
				Handler: b.RunFileCmd,
			},
		},
	}
}

func (b *Bot) RunCodeCmd(ctx sevcord.Ctx, args []any) {
	ctx.Modal(&sevcord.Modal{
		Title: "Run Code",
		Inputs: []sevcord.ModalInput{
			sevcord.ModalInput{
				Label:       "Code to Run",
				Style:       sevcord.ModalInputStyleParagraph,
				Placeholder: `[PRINT "Hello, World!"]`,
				Required:    true,
				MaxLength:   4000,
				MinLength:   1,
			},
		},
		Handler: func(ctx sevcord.Ctx, vals []string) {
			ctx.Acknowledge()

			// Actually run code
			d, err := b.Get(ctx.Guild())
			if err != nil {
				Error(ctx, err)
				return
			}
			err = b.RunCode("main.bsp", vals[0], ctx, NewExtensionCtx("_run", d, ctx))
			Error(ctx, err)
		},
	})
}

func (b *Bot) RunTagCmd(ctx sevcord.Ctx, vals []any) {
	dat, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	ctx.Acknowledge()

	// Get program
	prog, rsp := dat.GetProgram(vals[0].(string))
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}
	src, rsp := dat.GetSource(vals[0].(string))
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}

	// Run
	startTime := time.Now()
	err = b.RunCode(prog.ID+".bsp", src, ctx, NewExtensionCtx(vals[0].(string), dat, ctx))
	if Error(ctx, err) {
		return
	}

	// Increment uses
	prog.Uses++
	prog.LastUsed = startTime
	err = dat.SaveProgram(prog)
	Error(ctx, err)
}

func (b *Bot) RunFileCmd(ctx sevcord.Ctx, vals []any) {
	ctx.Acknowledge()

	resp, err := http.Get(vals[0].(*sevcord.SlashCommandAttachment).URL)
	if Error(ctx, err) {
		return
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if Error(ctx, err) {
		return
	}
	if len(dat) > 1048576 {
		ErrorMessage(ctx, "The maximum program size is **1MB**!")
		return
	}

	d, err := b.Get(ctx.Guild())
	if err != nil {
		Error(ctx, err)
		return
	}
	err = b.RunCode("main.bsp", string(dat), ctx, NewExtensionCtx("_run", d, ctx))
	Error(ctx, err)
}
