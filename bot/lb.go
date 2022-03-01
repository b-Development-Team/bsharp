package bot

import (
	"fmt"
	"sort"

	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/bwmarrin/discordgo"
)

const ItemsPerPage = 10

var btnRow = discordgo.ActionsRow{
	Components: []discordgo.MessageComponent{
		discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name:     "leftarrow",
				ID:       "861722690813165598",
				Animated: false,
			},
			CustomID: "prev",
		},
		discordgo.Button{
			Emoji: discordgo.ComponentEmoji{
				Name:     "rightarrow",
				ID:       "861722690926936084",
				Animated: false,
			},
			CustomID: "next",
		},
	},
}

func (b *Bot) LbUsedCmd(ctx *Ctx) {
	dat, err := b.Get(ctx.Guild())
	if ctx.Error(err) {
		return
	}

	ctx.Followup()

	// Load progs
	progs := make([]*db.Program, len(dat.Programs))
	dat.RLock()
	i := 0
	for _, p := range dat.Programs {
		progs[i] = p
		i++
	}
	dat.RUnlock()

	// Sort
	sort.Slice(progs, func(i, j int) bool {
		return progs[i].Uses > progs[j].Uses
	})

	// Calculate stats
	pages := len(progs) / ItemsPerPage
	if len(progs)%ItemsPerPage != 0 {
		pages++
	}
	page := 0

	// Send
	buildEmb := func() *discordgo.MessageEmbed {
		emb := &discordgo.MessageEmbed{
			Title:       "Most Used Tags",
			Color:       15844367, // Gold
			Description: "",
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page %d/%d", page+1, pages),
			},
		}
		for i := 0; i < ItemsPerPage; i++ {
			ind := i + page*ItemsPerPage
			if ind >= len(progs) {
				break
			}
			s := "s"
			uses := progs[ind].Uses
			if uses == 1 {
				s = ""
			}
			emb.Description += fmt.Sprintf("%d. **%s** - %d Use%s\n", i+1, progs[ind].Name, uses, s)
		}
		return emb
	}

	// Send
	ctx.Embed(buildEmb(), btnRow)
	ctx.BtnHandler(func(data discordgo.MessageComponentInteractionData, ctx *Ctx) {
		switch data.CustomID {
		case "prev":
			page--
			if page < 0 {
				page = pages - 1
			}
			ctx.Embed(buildEmb(), btnRow)
		case "next":
			page++
			if page >= pages {
				page = 0
			}
			ctx.Embed(buildEmb(), btnRow)
		}
	})
}
