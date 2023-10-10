package music

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/embed"
	"github.com/nabomhalang/halangcordgo/utils"
)

var (
	log *config.Logger = config.NewLogger("music-play")
)

type PlayCommandHandler struct{}

func (c *PlayCommandHandler) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "play",
		Description: "Plays a song from youtube or spotify playlist (or searches the query on youtube)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "link",
				Description: "Link or query to play",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "shuffle",
				Description: "Whether to shuffle the playlist or not",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "loop",
				Description: "Whether to loop the song or not",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "priority",
				Description: "Does this song have priority over the other songs in the queue?",
				Required:    false,
			},
		},
	}
}

func (c *PlayCommandHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if voiceChannel := utils.UserIsInVoiceChannel(s, i.GuildID, i.Member.User.ID); voiceChannel != nil {
		if utils.JoinVoiceChannel(s, i.Interaction, voiceChannel.ChannelID) {
			var (
				link                    string
				err                     error
				options                 = i.ApplicationCommandData().Options
				shuffle, loop, priority bool
			)

			for j := 1; j < len(options); j++ {
				switch options[j].Name {
				case "shuffle":
					shuffle = options[j].Value.(bool)
				case "loop":
					loop = options[j].Value.(bool)
				case "priority":
					priority = options[j].Value.(bool)
				}
			}

			link = options[0].Value.(string)

			if err == nil {
				play(s, link, i.Interaction, i.GuildID, i.Member.User.Username, shuffle, loop, priority)
			}
		}
	} else {
		embed.SendAndDeleteEmbedInteraction(s, embed.NewEmbed().SetTitle(s.State.User.Username).
			AddField("ERROR", fmt.Sprintf("User <@%s> is not in a voice channel", i.Member.User.ID), false).
			SetColor("red").MessageEmbed, i.Interaction, time.Second*10)
	}
}
