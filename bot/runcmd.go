package bot

import (
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) RunCodeCmd(ctx *Ctx) {
	err := ctx.Modal(&discordgo.InteractionResponseData{
		Title: "Run Code",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "code",
						Label:       "Code to Run",
						Style:       discordgo.TextInputParagraph,
						Placeholder: `[PRINT "Hello, World!"]`,
						Required:    true,
						MaxLength:   1536,
						MinLength:   1,
					},
				},
			},
		},
	}, func(dat discordgo.ModalSubmitInteractionData, ctx *Ctx) {
		ctx.Followup()

		// Actually run code
		src := dat.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		err := b.RunCode("main.bsp", src, ctx)
		ctx.Error(err)
	})
	ctx.Error(err)
}
