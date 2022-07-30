package bot

import (
	"fmt"
	"sort"

	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/Nv7-Github/sevcord"
)

const ItemsPerPage = 10

func LbCmd(b *Bot) sevcord.SlashCommandObject {
	return &sevcord.SlashCommand{
		Name:        "lb",
		Description: "View the most used tags!",
		Options:     []sevcord.Option{},
		Handler:     b.LbUsedCmd,
	}
}

func (b *Bot) LbUsedCmd(ctx sevcord.Ctx, args []any) {
	dat, err := b.Get(ctx.Guild())
	if Error(ctx, err) {
		return
	}

	ctx.Acknowledge()

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
	buildEmb := func() *sevcord.Response {
		desc := ""

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
			desc += fmt.Sprintf("%d. **%s** - %d Use%s\n", i+1, progs[ind].Name, uses, s)
		}

		return sevcord.EmbedResponse(sevcord.NewEmbedBuilder("Most Used Tags").
			Color(15844367).
			Footer(fmt.Sprintf("Page %d/%d", page+1, pages), "").
			Description(desc))
	}

	// Buttons
	var btns []sevcord.Component
	btns = []sevcord.Component{
		&sevcord.Button{
			Emoji: sevcord.ComponentEmojiCustom("leftarrow", "861722690813165598", false),
			Style: sevcord.ButtonStylePrimary,
			Handler: func(ctx sevcord.Ctx) {
				page--
				if page < 0 {
					page = pages - 1
				}
				ctx.Respond(buildEmb().ComponentRow(btns...))
			},
		},
		&sevcord.Button{
			Emoji: sevcord.ComponentEmojiCustom("rightarrow", "861722690926936084", false),
			Style: sevcord.ButtonStylePrimary,
			Handler: func(ctx sevcord.Ctx) {
				page++
				if page >= pages {
					page = 0
				}
				ctx.Respond(buildEmb().ComponentRow(btns...))
			},
		},
	}

	// Send
	ctx.Respond(buildEmb().ComponentRow(btns...))
}
