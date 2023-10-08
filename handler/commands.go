package handler

import (
	"github.com/bwmarrin/discordgo"
	Server "github.com/nabomhalang/halangcordgo/commands/server"
)

type CommandHandler interface {
	Command() *discordgo.ApplicationCommand
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CommandHandlers []CommandHandler = []CommandHandler{
	&Server.PingCommandHandler{},
}
