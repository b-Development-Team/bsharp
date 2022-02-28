package bot

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) CreateCmd(src string, id, name string, ctx *Ctx) {
	dat, err := b.Get(ctx.Guild())
	if ctx.Error(err) {
		return
	}
	ctx.Followup()

	// Create
	prog, rsp := dat.NewProgram(id, name, ctx.Author())
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
	err = dat.SaveProgram(prog)
	if ctx.Error(err) {
		return
	}
	err = dat.SaveSource(prog.ID, src)
	if ctx.Error(err) {
		return
	}
	ctx.Message(fmt.Sprintf("🆕 Created new program **%s**!", prog.Name))
}

func (b *Bot) CreateCodeCmd(ctx *Ctx) {
	err := ctx.Modal(&discordgo.InteractionResponseData{
		Title: "Create Program",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "id",
						Label:       "Program ID",
						Style:       discordgo.TextInputShort,
						Placeholder: `my-program`,
						Required:    true,
						MaxLength:   256,
						MinLength:   1,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "name",
						Label:       "Program Name",
						Style:       discordgo.TextInputShort,
						Placeholder: `My Program`,
						Required:    true,
						MaxLength:   256,
						MinLength:   1,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "code",
						Label:       "Program Code",
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
		id := dat.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		name := dat.Components[1].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		src := dat.Components[2].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		b.CreateCmd(src, id, name, ctx)
	})
	ctx.Error(err)
}

func (b *Bot) CreateFileCmd(id, name string, url string, ctx *Ctx) {
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

	b.CreateCmd(string(dat), id, name, ctx)
}
