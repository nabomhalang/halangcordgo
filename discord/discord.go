package discord

import (
	"strconv"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/commands"
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/server"
	"github.com/nabomhalang/halangcordgo/utils"
)

var (
	Server                = server.SV
	log    *config.Logger = config.NewLogger("discord")
)

func Ready(s *discordgo.Session, _ *discordgo.Ready) {
	err := s.UpdateGameStatus(0, "Serving "+strconv.Itoa(len(s.State.Guilds))+" guilds")
	if err != nil {
		log.Errorf("Failed to update game status: %s", err.Error())
	}
}

func GuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	utils.InitServer(g.ID)

	for _, c := range g.Channels {
		if c.Type == discordgo.ChannelTypeGuildVoice {
			Server[g.ID].VoiceChannelMembers[c.ID] = &atomic.Int32{}
		}
	}

	for _, v := range g.VoiceStates {
		Server[g.ID].VoiceChannelMembers[v.ChannelID].Add(1)
	}

	// log.Infof("Joined guild %s (%s) Guild %s has %d voice channels", g.Name, g.ID, g.Name, len(Server[g.ID].VoiceChannelMembers))
	Ready(s, nil)
}

func GuildDelete(s *discordgo.Session, g *discordgo.GuildDelete) {
	if Server[g.ID].IsPlaying() {
		utils.ClearAndExit(g.ID)
	}

	// log.Infof("Left guild %s (%s)", g.Name, g.ID)
	Ready(s, nil)
}

func VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v.UserID == s.State.User.ID && Server[v.GuildID].IsPlaying() {
		if v.ChannelID == "" {
			var err error

			Server[v.GuildID].VC, err = s.ChannelVoiceJoin(v.GuildID, Server[v.GuildID].VoiceChannel, false, true)
			if err != nil {
				log.Errorf("Failed to join voice channel: %s", err.Error())
				return
			}
		} else {
			Server[v.GuildID].VoiceChannel = v.ChannelID
		}
	}

	if v.ChannelID != "" {
		if v.BeforeUpdate != nil {
			Server[v.GuildID].VoiceChannelMembers[v.BeforeUpdate.ChannelID].Add(-1)
		}
		Server[v.GuildID].VoiceChannelMembers[v.ChannelID].Add(1)
	} else {
		Server[v.GuildID].VoiceChannelMembers[v.BeforeUpdate.ChannelID].Add(-1)
	}
}

func ChannelCreate(_ *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Type == discordgo.ChannelTypeGuildVoice && Server[c.GuildID].VoiceChannelMembers[c.ID] == nil {
		Server[c.GuildID].VoiceChannelMembers[c.ID] = &atomic.Int32{}
	}
}

func GuildMemberUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	if m.User.ID == s.State.User.ID && m.CommunicationDisabledUntil != nil && Server[m.GuildID].IsPlaying() {
		utils.ClearAndExit(m.GuildID)
	}
}

func RegisterCommands(s *discordgo.Session, commands map[string]commands.CommandHandler) {
	for _, handler := range commands {
		cmd := handler.Command()
		log.Infof("Registering command %s", cmd.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Errorf("Failed to register command %s: %s", cmd.Name, err.Error())
		}
	}
}
