package bot

import (
	"fmt"

	"github.com/Nv7-Github/sevcord"
)

func DescriptionCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommand{
		Name:        "description",
		Description: "Edit the description of a tag",
		Options: []sevcord.Option{
			{
				Name:         "tag",
				Description:  "The ID of the tag to edit",
				Required:     true,
				Autocomplete: b.Autocomplete,
				Kind:         sevcord.OptionKindString,
			},
		},
		Handler: b.DescriptionCmd,
	}
}

func (b *Bot) DescriptionCmd(ctx sevcord.Ctx, args []any) {
	ctx.Modal(&sevcord.Modal{
		Title: "Describe Tag",
		Inputs: []sevcord.ModalInput{
			{
				Label:       "New Description",
				Style:       sevcord.ModalInputStyleParagraph,
				Placeholder: `Description...`,
				Required:    true,
				MaxLength:   2048,
				MinLength:   1,
			},
		},
		Handler: func(ctx sevcord.Ctx, vals []string) {
			ctx.Acknowledge()

			// Actually run code
			data, err := b.Get(ctx.Guild())
			if Error(ctx, err) {
				return
			}
			prog, rsp := data.GetProgram(args[0].(string))
			if !rsp.Suc {
				ErrorMessage(ctx, rsp.Msg)
				return
			}
			if ctx.User().ID != prog.Creator {
				ErrorMessage(ctx, "Only the creator can change the description of their tag!")
				return
			}
			prog.Description = vals[0]
			err = data.SaveProgram(prog)
			if Error(ctx, err) {
				return
			}
			ctx.Respond(sevcord.MessageResponse(fmt.Sprintf("📝 Edited tag **%s**!", prog.Name)))
		},
	})
}

func ImageCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommand{
		Name:        "image",
		Description: "Edit the image of a tag!",
		Options: []sevcord.Option{
			{
				Name:         "tag",
				Description:  "The ID of the tag to edit!",
				Required:     true,
				Autocomplete: b.Autocomplete,
				Kind:         sevcord.OptionKindString,
			},
			{
				Name:        "image",
				Description: "The new image!",
				Required:    true,
				Kind:        sevcord.OptionKindAttachment,
			},
		},
		Handler: b.ImageCmd,
	}
}

func (b *Bot) ImageCmd(ctx sevcord.Ctx, args []any) {
	ctx.Acknowledge()

	data, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	prog, rsp := data.GetProgram(args[0].(string))
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}
	if ctx.User().ID != prog.Creator {
		ErrorMessage(ctx, "Only the creator can change the image of their tag!")
		return
	}
	prog.Image = args[1].(*sevcord.SlashCommandAttachment).URL
	err = data.SaveProgram(prog)
	if Error(ctx, err) {
		return
	}
	ctx.Respond(sevcord.MessageResponse(fmt.Sprintf("📝 Edited tag **%s**!", prog.Name)))
}
