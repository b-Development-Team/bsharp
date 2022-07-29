package bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Nv7-Github/sevcord"
)

func CreateCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommandGroup{
		Name:        "create",
		Description: "Create a tag!",
		Children: []sevcord.SlashCommandObject{
			&sevcord.SlashCommand{
				Name:        "tag",
				Description: "Create a tag from the source code of a B# program!",
				Options:     []sevcord.Option{},
				Handler:     b.CreateCodeCmd,
			},
			&sevcord.SlashCommand{
				Name:        "file",
				Description: "Create a tag from a file!",
				Options: []sevcord.Option{
					{
						Name:        "id",
						Description: "The ID of the tag to create!",
						Kind:        sevcord.OptionKindString,
						Required:    true,
					},
					{
						Name:        "name",
						Description: "The name of the tag to create!",
						Kind:        sevcord.OptionKindString,
						Required:    true,
					},
					{
						Name:        "file",
						Description: "The source code of the tag!",
						Kind:        sevcord.OptionKindAttachment,
						Required:    true,
					},
				},
				Handler: b.CreateFileCmd,
			},
		},
	}
}

func (b *Bot) CreateCmd(src string, id, name string, ctx sevcord.Ctx) {
	dat, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	ctx.Acknowledge()

	// Create
	prog, rsp := dat.NewProgram(id, name, ctx.User().ID)
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
	err = dat.SaveProgram(prog)
	if Error(ctx, err) {
		return
	}
	err = dat.SaveSource(prog.ID, src)
	if Error(ctx, err) {
		return
	}
	ctx.Respond(sevcord.MessageResponse(fmt.Sprintf("🆕 Created new tag **%s**!", prog.Name)))
}

func (b *Bot) CreateCodeCmd(ctx sevcord.Ctx, args []any) {
	ctx.Modal(&sevcord.Modal{
		Title: "Create Tag",
		Inputs: []sevcord.ModalInput{
			{
				Label:       "Program ID",
				Style:       sevcord.ModalInputStyleSentence,
				Placeholder: "my-program",
				Required:    true,
				MinLength:   1,
				MaxLength:   256,
			},
			{
				Label:       "Program Name",
				Style:       sevcord.ModalInputStyleSentence,
				Placeholder: "My Program",
				Required:    true,
				MinLength:   1,
				MaxLength:   256,
			},
			{
				Label:       "Source Code",
				Style:       sevcord.ModalInputStyleParagraph,
				Placeholder: `[PRINT "Hello, World!"]`,
				Required:    true,
				MinLength:   1,
				MaxLength:   1048576,
			},
		},
		Handler: func(ctx sevcord.Ctx, args []string) {
			b.CreateCmd(args[2], args[0], args[1], ctx)
		},
	})
}

func (b *Bot) CreateFileCmd(ctx sevcord.Ctx, args []any) {
	ctx.Acknowledge()

	resp, err := http.Get(args[2].(*sevcord.SlashCommandAttachment).URL)
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

	b.CreateCmd(string(dat), args[0].(string), args[1].(string), ctx)
}
