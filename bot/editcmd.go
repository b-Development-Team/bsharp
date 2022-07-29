package bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Nv7-Github/sevcord"
)

func EditCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommandGroup{
		Name:        "edit",
		Description: "Edit a B# tag!",
		Children: []sevcord.SlashCommandObject{
			&sevcord.SlashCommand{
				Name:        "tag",
				Description: "Edit a B# tag from the source code of a B# program!",
				Options: []sevcord.Option{
					{
						Name:         "id",
						Description:  "The ID of the tag to edit!",
						Required:     true,
						Autocomplete: b.Autocomplete,
						Kind:         sevcord.OptionKindString,
					},
				},
				Handler: b.EditCodeCmd,
			},
			&sevcord.SlashCommand{
				Name:        "file",
				Description: "Edit a B# tag from a file!",
				Options: []sevcord.Option{
					{
						Name:         "id",
						Description:  "The ID of the tag to edit!",
						Required:     true,
						Autocomplete: b.Autocomplete,
						Kind:         sevcord.OptionKindString,
					},
					{
						Name:        "file",
						Description: "The file to edit the tag from!",
						Required:    true,
						Kind:        sevcord.OptionKindAttachment,
					},
				},
				Handler: b.EditFileCmd,
			},
		},
	}
}

func (b *Bot) EditCmd(ctx sevcord.Ctx, src string, id string) {
	dat, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	ctx.Acknowledge()

	// Create
	prog, rsp := dat.GetProgram(id)
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}

	// Attempt to build
	_, err = b.BuildCode(prog.ID+".bsp", src, ctx)
	if Error(ctx, err) {
		return
	}

	// Save
	err = dat.SaveSource(prog.ID, src)
	if Error(ctx, err) {
		return
	}
	ctx.Respond(sevcord.MessageResponse(fmt.Sprintf("📝 Edited tag **%s**!", prog.Name)))
}

func (b *Bot) EditCodeCmd(ctx sevcord.Ctx, args []any) {
	ctx.Modal(&sevcord.Modal{
		Title: "Edit Code",
		Inputs: []sevcord.ModalInput{
			{
				Label:       "New Code",
				Style:       sevcord.ModalInputStyleParagraph,
				Placeholder: `[PRINT "Hello, World!"]`,
				Required:    true,
				MaxLength:   4000,
				MinLength:   1,
			},
		},
		Handler: func(ctx sevcord.Ctx, vals []string) {
			ctx.Acknowledge()

			b.EditCmd(ctx, vals[0], args[0].(string))
		},
	})
}

func (b *Bot) EditFileCmd(ctx sevcord.Ctx, args []any) {
	ctx.Acknowledge()

	resp, err := http.Get(args[1].(*sevcord.SlashCommandAttachment).URL)
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

	b.EditCmd(ctx, string(dat), args[0].(string))
}
