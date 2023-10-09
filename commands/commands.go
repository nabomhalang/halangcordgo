package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/config"
)

type Command struct {
	CommandHandlers map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var log = config.NewLogger("commands")
