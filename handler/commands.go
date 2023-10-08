package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/commands/server"
)

type CommandHandler interface {
	Command() *discordgo.ApplicationCommand
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CommandHandlers []CommandHandler = []CommandHandler{
	&server.PingCommandHandler{},
}
