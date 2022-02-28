package bot

import (
	"time"

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

func (b *Bot) RunTagCmd(id string, ctx *Ctx) {
	dat, err := b.Get(ctx.Guild())
	if ctx.Error(err) {
		return
	}
	ctx.Followup()

	// Get program
	prog, rsp := dat.GetProgram(id)
	if !rsp.Suc {
		ctx.ErrorMessage(rsp.Msg)
		return
	}
	src, rsp := dat.GetSource(id)
	if !rsp.Suc {
		ctx.ErrorMessage(rsp.Msg)
		return
	}

	// Run
	startTime := time.Now()
	err = b.RunCode(prog.ID+".bsp", src, ctx)
	if ctx.Error(err) {
		return
	}

	// Increment uses
	prog.Uses++
	prog.LastUsed = startTime
	err = dat.SaveProgram(prog)
	ctx.Error(err)
}
