package embed

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func SendAndDeleteEmbedInteraction(s *discordgo.Session, embed *discordgo.MessageEmbed, i *discordgo.Interaction, wait time.Duration) {
	SendEmbedInteraction(s, embed, i, nil)

	time.Sleep(wait)

	err := s.InteractionResponseDelete(i)
	if err != nil {
		log.Errorf("Failed to delete interaction response: %s", err.Error())
	}
}

func SendEmbedInteraction(s *discordgo.Session, embed *discordgo.MessageEmbed, i *discordgo.Interaction, c chan<- struct{}) {
	sliceEmbed := []*discordgo.MessageEmbed{embed}

	err := s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: sliceEmbed,
		},
	})
	if err != nil {
		log.Errorf("Failed to send interaction response: %s", err.Error())
		return
	}

	if c != nil {
		c <- struct{}{}
	}
}
