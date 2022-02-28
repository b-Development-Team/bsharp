package bot

import (
	"sync"

	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	*sync.RWMutex
	*db.DB

	dg     *discordgo.Session // use Ctx.DG whenever possible
	modals map[string]func(discordgo.ModalSubmitInteractionData, *Ctx)
	appID  string
}

func NewBot(path string, token string, appID string, guild string) (*Bot, error) {
	b := &Bot{
		RWMutex: &sync.RWMutex{},
		appID:   appID,
		modals:  make(map[string]func(discordgo.ModalSubmitInteractionData, *Ctx)),
	}
	d, err := db.NewDB(path)
	if err != nil {
		return nil, err
	}
	b.DB = d
	err = b.initDG(token, appID, guild)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Bot) Close() {
	b.DB.Close()
	b.dg.Close()
}
