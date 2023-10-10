package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/commands/music"
	"github.com/nabomhalang/halangcordgo/commands/server"
	"github.com/nabomhalang/halangcordgo/config"
)

type CommandHandler interface {
	Command() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	log             *config.Logger            = config.NewLogger("commands")
	CommandHandlers map[string]CommandHandler = map[string]CommandHandler{
		// Server commands
		"ping":    &server.PingCommandHandler{},
		"restart": &server.RestartCommandHandler{},

		// Music commands
		"play": &music.PlayCommandHandler{},
	}
)
