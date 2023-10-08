package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/handler"
)

func Register(s *discordgo.Session, commands []handler.CommandHandler) {
	for _, handler := range commands {
		cmd := handler.Command()
		log.Infof("Registering command: %s", cmd.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Errorf("Failed to register command: %s", cmd.Name)
			continue
		}
	}

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			for _, handler := range commands {
				if i.ApplicationCommandData().Name == handler.Command().Name {
					handler.Handle(s, i)
					break
				}
			}
		}
	})
}
