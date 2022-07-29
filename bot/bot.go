package bot

import (
	"sync"

	"github.com/Nv7-Github/bsharp/bot/db"
	"github.com/Nv7-Github/sevcord"
)

type Bot struct {
	*sync.RWMutex
	*db.DB

	c *sevcord.Client
}

func ErrorMessage(ctx sevcord.Ctx, msg string) {
	ctx.Respond(sevcord.EmbedResponse(sevcord.NewEmbedBuilder("Error").Color(15548997).Description(msg)))
}

func Error(ctx sevcord.Ctx, err error) bool {
	if err != nil {
		ErrorMessage(ctx, err.Error())
		return true
	}
	return false
}

func (b *Bot) Autocomplete(ctx sevcord.Ctx, val any) []sevcord.Choice {
	db, err := b.Get(ctx.Guild())
	if err != nil {
		return nil
	}
	return db.Autocomplete(val.(string))
}

func NewBot(path string, token string) (*Bot, error) {
	b := &Bot{
		RWMutex: &sync.RWMutex{},
	}
	d, err := db.NewDB(path)
	if err != nil {
		return nil, err
	}
	b.DB = d
	c, err := sevcord.NewClient(token)
	if err != nil {
		return nil, err
	}
	c.HandleSlashCommand(BuildCmd(b))
	c.Start()
	return b, nil
}

func (b *Bot) Close() {
	b.DB.Close()
	b.c.Close()
}
