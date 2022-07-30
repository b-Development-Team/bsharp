package bot

import (
	"fmt"
	"strings"

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

func InfoCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommand{
		Name:        "info",
		Description: "Get info about a tag!",
		Options: []sevcord.Option{
			{
				Name:         "tag",
				Description:  "The ID of the tag to get info about!",
				Required:     true,
				Autocomplete: b.Autocomplete,
				Kind:         sevcord.OptionKindString,
			},
		},
		Handler: b.InfoCmd,
	}
}

func (b *Bot) InfoCmd(ctx sevcord.Ctx, args []any) {
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

	emb := sevcord.NewEmbedBuilder(fmt.Sprintf("%s Info", prog.Name)).
		Description(prog.Description).
		Color(3447003).
		Field("ID", prog.ID, false).
		Field("Uses", fmt.Sprintf("%d", prog.Uses), false).
		Field("Creator", fmt.Sprintf("<@%s>", prog.Creator), false).
		Field("Created", fmt.Sprintf("<t:%d>", prog.Created.Unix()), false).
		Field("Last Used", fmt.Sprintf("<t:%d>", prog.LastUsed.Unix()), false)

	if prog.Image != "" {
		emb.Thumbnail(prog.Image)
	}

	ctx.Respond(sevcord.EmbedResponse(emb))
}

func SourceCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommand{
		Name:        "source",
		Description: "Get the source code of a tag!",
		Options: []sevcord.Option{
			{
				Name:         "tag",
				Description:  "The ID of the tag to get info about!",
				Required:     true,
				Autocomplete: b.Autocomplete,
				Kind:         sevcord.OptionKindString,
			},
		},
		Handler: b.SourceCmd,
	}
}

func (b *Bot) SourceCmd(ctx sevcord.Ctx, args []any) {
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
	src, rsp := data.GetSource(args[0].(string))
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}

	ctx.Respond(sevcord.MessageResponse(fmt.Sprintf("Source for **%s**:", prog.Name)).File(sevcord.File{
		Name:        prog.ID + ".txt",
		ContentType: "text/plain",
		Reader:      strings.NewReader(src),
	}))
}
