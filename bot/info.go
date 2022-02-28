package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) InfoCmd(id string, ctx *Ctx) {
	ctx.Followup()

	data, err := b.Get(ctx.Guild())
	if ctx.Error(err) {
		return
	}
	prog, rsp := data.GetProgram(id)
	if !rsp.Suc {
		ctx.ErrorMessage(rsp.Msg)
		return
	}

	var thumbnail *discordgo.MessageEmbedThumbnail
	if prog.Image != "" {
		thumbnail = &discordgo.MessageEmbedThumbnail{URL: prog.Image}
	}

	ctx.Embed(&discordgo.MessageEmbed{
		Title:       fmt.Sprintf("%s Info", prog.Name),
		Description: prog.Description,
		Color:       3447003, // Blue
		Thumbnail:   thumbnail,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "ID",
				Value: prog.ID,
			},
			{
				Name:  "Uses",
				Value: fmt.Sprintf("%d", prog.Uses),
			},
			{
				Name:  "Creator",
				Value: fmt.Sprintf("<@%s>", prog.Creator),
			},
			{
				Name:  "Created",
				Value: fmt.Sprintf("<t:%d>", prog.Created.Unix()),
			},
			{
				Name:  "Last Used",
				Value: fmt.Sprintf("<t:%d>", prog.LastUsed.Unix()),
			},
		},
	})
}

func (b *Bot) SourceCmd(id string, ctx *Ctx) {
	ctx.Followup()

	data, err := b.Get(ctx.Guild())
	if ctx.Error(err) {
		return
	}
	prog, rsp := data.GetProgram(id)
	if !rsp.Suc {
		ctx.ErrorMessage(rsp.Msg)
		return
	}
	src, rsp := data.GetSource(id)
	if !rsp.Suc {
		ctx.ErrorMessage(rsp.Msg)
		return
	}

	ctx.Message(fmt.Sprintf("Source for **%s**:", prog.Name), &discordgo.File{
		Name:        prog.ID + ".txt",
		ContentType: "text/plain",
		Reader:      strings.NewReader(src),
	})
}
