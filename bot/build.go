package bot

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) BuildCodeCmd(ctx *Ctx) {
	err := ctx.Modal(&discordgo.InteractionResponseData{
		Title: "Build Code",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "code",
						Label:       "Code to Build",
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

		// Actually run code
		src := dat.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
		code, err := b.CompileCode("main.bsp", src, ctx)
		if ctx.Error(err) {
			return
		}
		ctx.Message("", &discordgo.File{
			Name: "main.c",
			//ContentType: "text/plain",
			Reader: strings.NewReader(code),
		})
	})
	ctx.Error(err)
}

func (b *Bot) BuildTagCmd(id string, ctx *Ctx) {
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
	code, err := b.CompileCode(prog.ID+".bsp", src, ctx)
	if ctx.Error(err) {
		return
	}
	ctx.Message(fmt.Sprintf("🖨️ Compiled code for **%s**", prog.Name), &discordgo.File{
		Name: "main.c",
		//ContentType: "text/plain",
		Reader: strings.NewReader(code),
	})
}

func (b *Bot) BuildFileCmd(url string, ctx *Ctx) {
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

	code, err := b.CompileCode("main.bsp", string(dat), ctx)
	if ctx.Error(err) {
		return
	}
	ctx.Message("", &discordgo.File{
		Name: "main.c",
		//ContentType: "text/plain",
		Reader: strings.NewReader(code),
	})
}
