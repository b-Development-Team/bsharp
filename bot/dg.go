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
	isButton     bool
	hasResponded *string
}

func (c *Ctx) Guild() string {
	return c.i.GuildID
}

func (c *Ctx) Author() string {
	return c.i.Member.User.ID
}

func (c *Ctx) BtnHandler(handler func(data discordgo.MessageComponentInteractionData, ctx *Ctx)) {
	c.b.Lock()
	c.b.btns[c.i.Message.ID] = handler
	c.b.Unlock()
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

func (c *Ctx) Message(m string, attachments ...*discordgo.File) error {
	c.Lock()
	defer c.Unlock()

	if c.followup {
		if c.hasResponded != nil {
			_, err := c.DG.FollowupMessageEdit(c.i, *c.hasResponded, &discordgo.WebhookEdit{
				Content: m,
				Embeds:  make([]*discordgo.MessageEmbed, 0),
				Files:   attachments,
			})
			return err
		}
		msg, err := c.DG.FollowupMessageCreate(c.i, true, &discordgo.WebhookParams{
			Content: m,
			Files:   attachments,
		})
		c.i.Message = msg
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
			Files:   attachments,
		},
	})
}

func (c *Ctx) Embed(emb *discordgo.MessageEmbed, components ...discordgo.MessageComponent) error {
	c.Lock()
	defer c.Unlock()

	if c.isButton {
		err := c.DG.InteractionRespond(c.i, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content:    "",
				Embeds:     []*discordgo.MessageEmbed{emb},
				Components: components,
			},
		})
		c.isButton = false
		return err
	}

	if c.followup {
		if c.hasResponded != nil {
			_, err := c.DG.FollowupMessageEdit(c.i, *c.hasResponded, &discordgo.WebhookEdit{
				Content:    "",
				Embeds:     []*discordgo.MessageEmbed{emb},
				Components: components,
			})
			return err
		}
		msg, err := c.DG.FollowupMessageCreate(c.i, true, &discordgo.WebhookParams{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
		})
		c.i.Message = msg
		if err != nil {
			return err
		}
		c.hasResponded = &msg.ID
		return nil
	}
	return c.DG.InteractionRespond(c.i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{emb},
			Components: components,
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
	b.debug = guild != ""
	dg, err := discordgo.New("Bot " + strings.TrimSpace(token))
	if err != nil {
		return err
	}
	b.dg = dg
	if b.debug {
		for i, cmd := range commands {
			commands[i].Name = "test" + cmd.Name
		}
	}
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
		if b.debug {
			d.Name = strings.TrimPrefix(d.Name, "test")
		}
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

	case discordgo.InteractionMessageComponent:
		d := i.MessageComponentData()
		b.RLock()
		h, ok := b.btns[i.Message.ID]
		b.RUnlock()
		if ok {
			ctx := &Ctx{
				i:        i.Interaction,
				DG:       s,
				b:        b,
				isButton: true,
			}
			h(d, ctx)
		}

	case discordgo.InteractionApplicationCommandAutocomplete:
		d := i.ApplicationCommandData()
		if b.debug {
			if !strings.HasPrefix(d.Name, "test") {
				return
			}
			d.Name = strings.TrimPrefix(d.Name, "test")
		}
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
