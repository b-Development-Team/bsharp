package bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) EditCmd(src string, id string, ctx *Ctx) {
	dat, err := b.Get(ctx.Guild())
	if ctx.Error(err) {
		return
	}
	ctx.Followup()

	// Create
	prog, rsp := dat.GetProgram(id)
	if !rsp.Suc {
		ctx.ErrorMessage(rsp.Msg)
		return
	}

	// Attempt to build
	_, err = b.BuildCode(prog.ID+".bsp", src, ctx)
	if ctx.Error(err) {
		return
	}

	// Save
	err = dat.SaveSource(prog.ID, src)
	if ctx.Error(err) {
		return
	}
	ctx.Message(fmt.Sprintf("🖋️ Edited tag **%s**!", prog.Name))
}

func (b *Bot) EditCodeCmd(id string, ctx *Ctx) {
	err := ctx.Modal(&discordgo.InteractionResponseData{
		Title: "Edit Code",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "code",
						Label:       "New Code",
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

		src := dat.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		b.EditCmd(src, id, ctx)
	})
	ctx.Error(err)
}

func (b *Bot) EditFileCmd(id, url string, ctx *Ctx) {
	ctx.Followup()

	resp, err := http.Get(url)
	if ctx.Error(err) {
		return
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if ctx.Error(err) {
		return
	}
	if len(dat) > 1048576 {
		ctx.ErrorMessage("The maximum program size is **1MB**!")
		return
	}

	b.EditCmd(string(dat), id, ctx)
}
