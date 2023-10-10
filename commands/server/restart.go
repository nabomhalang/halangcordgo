package server

import (
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/embed"
)

var (
	owner string = config.Get().Owner
)

type RestartCommandHandler struct{}

func (c *RestartCommandHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "restart",
		Description: "Restart the bot",
	}
}

func (c *RestartCommandHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if owner == i.Member.User.ID {
		embed.SendAndDeleteEmbedInteraction(s, embed.NewEmbed().SetTitle(s.State.User.Username).AddField("Restart", "Bot restarting now...", false).
			SetColor(0x7289DA).MessageEmbed, i.Interaction, time.Second*1)
		_ = s.Close()
		os.Exit(0)
	} else {
		embed.SendAndDeleteEmbedInteraction(s,
			embed.NewEmbed().
				SetTitle(s.State.User.Username).
				AddField("Error", fmt.Sprintf("I'm sorry <@%s>, I'm afraid I can't do that", i.Member.User.ID), false).
				SetColor("red").
				MessageEmbed, i.Interaction, time.Second*5)
	}
}
