package bot

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "run",
			Description: "Runs a B# program!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "code",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Run the source code of a B# program!",
					Options:     []*discordgo.ApplicationCommandOption{}, // Code is entered through a modal
				},
			},
		},
	}

	handlers = map[string]func(*Ctx, discordgo.ApplicationCommandInteractionData, *Bot){
		"run": func(ctx *Ctx, dat discordgo.ApplicationCommandInteractionData, b *Bot) {
			data := dat.Options[0]
			switch data.Name {
			case "code":
				b.RunCodeCmd(ctx)
			}
		},
	}
)
