package bot

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Nv7-Github/sevcord"
)

func BuildCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommandGroup{
		Name:        "build",
		Description: "Compile a B# program!",
		Children: []sevcord.SlashCommandObject{
			&sevcord.SlashCommand{
				Name:        "code",
				Description: "Build the source code of a B# program!",
				Options:     []sevcord.Option{},
				Handler:     b.BuildCodeCmd,
			},
			&sevcord.SlashCommand{
				Name:        "tag",
				Description: "Build a tag!",
				Options: []sevcord.Option{
					{
						Name:         "tag",
						Description:  "The ID of the tag to build.",
						Required:     true,
						Autocomplete: b.Autocomplete,
						Kind:         sevcord.OptionKindString,
					},
				},
				Handler: b.BuildTagCmd,
			},
			&sevcord.SlashCommand{
				Name:        "file",
				Description: "Build the source code of a B# program, uploaded as a file!!",
				Options: []sevcord.Option{
					{
						Name:        "file",
						Description: "The file to build.",
						Required:    true,
						Kind:        sevcord.OptionKindAttachment,
					},
				},
				Handler: b.BuildTagCmd,
			},
		},
	}
}

func (b *Bot) BuildCodeCmd(ctx sevcord.Ctx, args []any) {
	ctx.Modal(&sevcord.Modal{
		Title: "Build Code",
		Inputs: []sevcord.ModalInput{
			{
				Label:       "Code to Build",
				Style:       sevcord.ModalInputStyleParagraph,
				Placeholder: `[PRINT "Hello, World!"]`,
				Required:    true,
				MinLength:   1,
				MaxLength:   4000,
			},
		},
		Handler: func(ctx sevcord.Ctx, vals []string) {
			ctx.Acknowledge()

			// Actually run code
			code, err := b.CompileCode("main.bsp", vals[0], ctx)
			if Error(ctx, err) {
				return
			}
			ctx.Respond(sevcord.MessageResponse("").File(sevcord.File{
				Name: "main.c",
				//ContentType: "text/plain",
				Reader: strings.NewReader(code),
			}))
		},
	})
}

func (b *Bot) BuildTagCmd(ctx sevcord.Ctx, args []any) {
	dat, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	ctx.Acknowledge()

	// Get program
	id := args[0].(string)
	prog, rsp := dat.GetProgram(id)
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}
	src, rsp := dat.GetSource(id)
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}

	// Run
	code, err := b.CompileCode(prog.ID+".bsp", src, ctx)
	if Error(ctx, err) {
		return
	}
	ctx.Respond(sevcord.MessageResponse(fmt.Sprintf("🖨️ Compiled code for **%s**", prog.Name)).File(sevcord.File{
		Name: "main.c",
		//ContentType: "text/plain",
		Reader: strings.NewReader(code),
	}))
}

func (b *Bot) BuildFileCmd(ctx sevcord.Ctx, args []any) {
	ctx.Acknowledge()

	resp, err := http.Get(args[0].(*sevcord.SlashCommandAttachment).URL)
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

	code, err := b.CompileCode("main.bsp", string(dat), ctx)
	if Error(ctx, err) {
		return
	}
	ctx.Respond(sevcord.MessageResponse("").File(sevcord.File{
		Name: "main.c",
		//ContentType: "text/plain",
		Reader: strings.NewReader(code),
	}))
}
