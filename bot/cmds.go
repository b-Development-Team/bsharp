package bot

import (
	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/bwmarrin/discordgo"
)

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
				{
					Name:        "tag",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Run a tag!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:         discordgo.ApplicationCommandOptionString,
							Name:         "id",
							Description:  "The ID of the tag to run.",
							Required:     true,
							Autocomplete: true,
						},
					}, // Code is entered through a modal
				},
			},
		},
		{
			Name:        "create",
			Description: "Creates a B# tag!",
			Type:        discordgo.ChatApplicationCommand,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "tag",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Create a tag from the source code of a B# program!",
					Options:     []*discordgo.ApplicationCommandOption{},
				},
				{
					Name:        "file",
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Description: "Create a tag from a file!",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "id",
							Description: "The ID of the tag to create!",
							Required:    true,
						},

						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "The name of the tag to create!",
							Required:    true,
						},
						{
							Type:        discordgo.ApplicationCommandOptionAttachment,
							Name:        "file",
							Description: "The file to create a tag from!",
							Required:    true,
						},
					}, // Code is entered through a modal
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

			case "tag":
				b.RunTagCmd(data.Options[0].StringValue(), ctx)
			}
		},
		"create": func(ctx *Ctx, dat discordgo.ApplicationCommandInteractionData, b *Bot) {
			data := dat.Options[0]
			switch data.Name {
			case "tag":
				b.CreateCodeCmd(ctx)

			case "file":
				var url string
				var name string
				var id string
				for _, opt := range data.Options {
					switch opt.Name {
					case "id":
						id = opt.StringValue()

					case "name":
						name = opt.StringValue()

					case "file":
						id := opt.Value.(string)
						attachment := dat.Resolved.Attachments[id]
						url = attachment.URL
					}
				}

				b.CreateFileCmd(id, name, url, ctx)
			}
		},
	}

	autocomplete = map[string]func(*db.Data, discordgo.ApplicationCommandInteractionData) []*discordgo.ApplicationCommandOptionChoice{
		"run": func(d *db.Data, dat discordgo.ApplicationCommandInteractionData) []*discordgo.ApplicationCommandOptionChoice {
			return d.Autocomplete(dat.Options[0].Options[0].StringValue())
		},
	}
)
