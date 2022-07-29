package bot

import (
	"io"
	"net/http"
	"time"

	"github.com/Nv7-Github/sevcord"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) RunCodeCmd(ctx sevcord.Ctx) {
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
						MaxLength:   4000,
						MinLength:   1,
					},
				},
			},
		},
	}, func(dat discordgo.ModalSubmitInteractionData, ctx sevcord.Ctx) {
		ctx.Acknowledge()

		// Actually run code
		src := dat.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		d, err := b.Get(ctx.Guild())
		if err != nil {
			Error(ctx, err)
			return
		}
		err = b.RunCode("main.bsp", src, ctx, NewExtensionCtx("_run", d, ctx))
		Error(ctx, err)
	})
	Error(ctx, err)
}

func (b *Bot) RunTagCmd(id string, ctx sevcord.Ctx) {
	dat, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}
	ctx.Acknowledge()

	// Get program
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
	startTime := time.Now()
	err = b.RunCode(prog.ID+".bsp", src, ctx, NewExtensionCtx(id, dat, ctx))
	if Error(ctx, err) {
		return
	}

	// Increment uses
	prog.Uses++
	prog.LastUsed = startTime
	err = dat.SaveProgram(prog)
	Error(ctx, err)
}

func (b *Bot) RunFileCmd(url string, ctx sevcord.Ctx) {
	ctx.Acknowledge()

	resp, err := http.Get(url)
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
