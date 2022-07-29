package bot

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/sevcord"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) InfoCmd(id string, ctx sevcord.Ctx) {
	ctx.Acknowledge()

	data, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	prog, rsp := data.GetProgram(id)
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
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

func (b *Bot) SourceCmd(id string, ctx sevcord.Ctx) {
	ctx.Acknowledge()

	data, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	prog, rsp := data.GetProgram(id)
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}
	src, rsp := data.GetSource(id)
	if !rsp.Suc {
		ErrorMessage(ctx, rsp.Msg)
		return
	}

	ctx.Message(fmt.Sprintf("Source for **%s**:", prog.Name), &discordgo.File{
		Name:        prog.ID + ".txt",
		ContentType: "text/plain",
		Reader:      strings.NewReader(src),
	})
}
