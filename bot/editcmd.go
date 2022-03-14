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
	ctx.Message(fmt.Sprintf("📝 Edited tag **%s**!", prog.Name))
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
						MaxLength:   4000,
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

func (b *Bot) DescriptionCmd(id string, ctx *Ctx) {
	err := ctx.Modal(&discordgo.InteractionResponseData{
		Title: "Describe Tag",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "description",
						Label:       "New Description",
						Style:       discordgo.TextInputParagraph,
						Placeholder: `Description...`,
						Required:    true,
						MaxLength:   2048,
						MinLength:   1,
					},
				},
			},
		},
	}, func(dat discordgo.ModalSubmitInteractionData, ctx *Ctx) {
		ctx.Followup()

		// Actually run code
		desc := dat.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		data, err := b.Get(ctx.Guild())
		if ctx.Error(err) {
			return
		}
		prog, rsp := data.GetProgram(id)
		if !rsp.Suc {
			ctx.ErrorMessage(rsp.Msg)
			return
		}
		if ctx.Author() != prog.Creator {
			ctx.ErrorMessage("Only the creator can change the description of their tag!")
			return
		}
		prog.Description = desc
		err = data.SaveProgram(prog)
		if ctx.Error(err) {
			return
		}
		ctx.Message(fmt.Sprintf("📝 Edited tag **%s**!", prog.Name))
	})
	ctx.Error(err)
}

func (b *Bot) ImageCmd(id, url string, ctx *Ctx) {
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
	if ctx.Author() != prog.Creator {
		ctx.ErrorMessage("Only the creator can change the image of their tag!")
		return
	}
	prog.Image = url
	err = data.SaveProgram(prog)
	if ctx.Error(err) {
		return
	}
	ctx.Message(fmt.Sprintf("📝 Edited tag **%s**!", prog.Name))
}
