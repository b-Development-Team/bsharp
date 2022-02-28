package bot

import (
	"strings"
	"sync"

	"github.com/bwmarrin/discordgo"
)

type Ctx struct {
	sync.Mutex // make sure everything happens in order

	i  *discordgo.Interaction
	DG *discordgo.Session
	b  *Bot

	followup     bool
	hasResponded *string
}

func (c *Ctx) Guild() string {
	return c.i.GuildID
}

func (c *Ctx) Author() string {
	return c.i.Member.User.ID
}

func (c *Ctx) Resp(m string) error {
	c.Lock()
	defer c.Unlock()

	if c.followup {
		return c.Message(m)
	}
	return c.DG.InteractionRespond(c.i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
			Flags:   1 << 6,
		},
	})
}

func (c *Ctx) Message(m string) error {
	c.Lock()
	defer c.Unlock()

	if c.followup {
		if c.hasResponded != nil {
			_, err := c.DG.FollowupMessageEdit(c.b.appID, c.i, *c.hasResponded, &discordgo.WebhookEdit{
				Content: m,
				Embeds:  make([]*discordgo.MessageEmbed, 0),
			})
			return err
		}
		msg, err := c.DG.FollowupMessageCreate(c.b.appID, c.i, true, &discordgo.WebhookParams{
			Content: m,
		})
		if err != nil {
			return err
		}
		c.hasResponded = &msg.ID
		return nil
	}
	return c.DG.InteractionRespond(c.i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: m,
		},
	})
}

func (c *Ctx) Embed(emb *discordgo.MessageEmbed) error {
	c.Lock()
	defer c.Unlock()

	if c.followup {
		if c.hasResponded != nil {
			_, err := c.DG.FollowupMessageEdit(c.b.appID, c.i, *c.hasResponded, &discordgo.WebhookEdit{
				Content: "",
				Embeds:  []*discordgo.MessageEmbed{emb},
			})
			return err
		}
		msg, err := c.DG.FollowupMessageCreate(c.b.appID, c.i, true, &discordgo.WebhookParams{
			Embeds: []*discordgo.MessageEmbed{emb},
		})
		if err != nil {
			return err
		}
		c.hasResponded = &msg.ID
		return nil
	}
	return c.DG.InteractionRespond(c.i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{emb},
		},
	})
}

func (c *Ctx) Followup() error {
	c.Lock()
	defer c.Unlock()

	if c.followup {
		return nil
	}
	c.followup = true
	return c.DG.InteractionRespond(c.i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func (c *Ctx) ErrorMessage(err string) error {
	return c.Embed(&discordgo.MessageEmbed{
		Title:       "Error",
		Color:       15548997, // Red
		Description: err,
	})
}

func (c *Ctx) Error(err error) bool {
	if err != nil {
		c.ErrorMessage(err.Error())
		return true
	}
	return false
}

func (c *Ctx) Modal(m *discordgo.InteractionResponseData, handler func(discordgo.ModalSubmitInteractionData, *Ctx)) error {
	m.CustomID = c.i.ID
	err := c.b.dg.InteractionRespond(c.i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: m,
	})
	if err != nil {
		return err
	}
	c.b.Lock()
	c.b.modals[c.i.ID] = handler
	c.b.Unlock()
	return nil
}

func (b *Bot) initDG(token string, appID string, guild string) error {
	dg, err := discordgo.New("Bot " + strings.TrimSpace(token))
	if err != nil {
		return err
	}
	b.dg = dg
	_, err = b.dg.ApplicationCommandBulkOverwrite(appID, guild, commands)
	if err != nil {
		return err
	}
	dg.AddHandler(b.InteractionHandler)

	return b.dg.Open()
}

func (b *Bot) InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ctx := &Ctx{
		i:  i.Interaction,
		DG: s,
		b:  b,
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		d := i.ApplicationCommandData()
		h, ok := handlers[d.Name]
		if ok {
			h(ctx, d, b)
		}

	case discordgo.InteractionModalSubmit:
		d := i.ModalSubmitData()
		b.RLock()
		h, ok := b.modals[d.CustomID]
		b.RUnlock()
		if ok {
			ctx := &Ctx{
				i:  i.Interaction,
				DG: s,
				b:  b,
			}
			h(d, ctx)
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		d := i.ApplicationCommandData()
		h, ok := autocomplete[d.Name]
		if ok {
			dat, err := b.Get(i.GuildID)
			if err != nil {
				return
			}
			res := h(dat, d)
			if res != nil {
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionApplicationCommandAutocompleteResult,
					Data: &discordgo.InteractionResponseData{
						Choices: res,
					},
				})
			}
		}
	}
}
