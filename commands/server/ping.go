package server

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/embed"
)

type PingCommandHandler struct{}

func (h *PingCommandHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "You say ping, I say pong",
	}
}

func (c *PingCommandHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed.SendEmbedInteraction(
		s,
		embed.NewEmbed().SetTitle("Pong!").SetDescription("Ping Pong message").SetColor("random").MessageEmbed,
		i.Interaction,
		nil,
	)
}
