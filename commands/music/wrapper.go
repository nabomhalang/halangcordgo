package music

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/embed"
	"github.com/nabomhalang/halangcordgo/utils"
)

func play(s *discordgo.Session, song string, i *discordgo.Interaction, guild, username string, random, loop, priority bool) {
	switch {
	// case utils.IsValidURL(song):
	// 	downloadAndPlay(s, guild, song, username, i, random, loop, true, priority)
	default:
		video, err := utils.SearchSong(song)
		if err == nil {
			utils.DownloadAndPlay(s, guild, username, video, i, random, loop, true, priority)
		} else {
			embed.SendAndDeleteEmbedInteraction(s,
				embed.NewEmbed().
					SetTitle(fmt.Sprintf("<@%s>", s.State.User.ID)).
					AddField("ERROR", err.Error(), false).
					SetColor("red").
					MessageEmbed, i, time.Second*10)
		}
	}
}
