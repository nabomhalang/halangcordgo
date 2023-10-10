package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/commands"
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/discord"
	"github.com/nabomhalang/halangcordgo/embed"
	"github.com/nabomhalang/halangcordgo/utils"
	"github.com/nabomhalang/halangcordgo/youtube"
)

var (
	Server           = discord.Server
	Owner     string = config.Get().Owner
	Token     string = config.Get().Token
	BlackList map[string]bool
	s         *discordgo.Session
	Youtube   *youtube.Youtube = youtube.YT
	log       *config.Logger   = config.NewLogger("main")
)

func init() {
	var err error

	if _, err = os.Stat(config.Get().CachePath); err != nil {
		if err = os.Mkdir(config.Get().CachePath, 0755); err != nil {
			log.Errorf("Failed to create cache directory: %s", err.Error())
		}
	}

	os.Remove("--Frag1")

	if utils.IsCommandNotAvailable("dca") {
		log.Fatal("dca is not installed")
	}

	if utils.IsCommandNotAvailable("ffmpeg") {
		log.Fatal("ffmpeg is not installed")
	}

	if utils.IsCommandNotAvailable("yt-dlp") {
		log.Fatal("yt-dlp is not installed")
	}

}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Failed to create discord session: %s", err.Error())
		return
	}
	defer dg.Close()

	err = dg.Open()
	if err != nil {
		log.Errorf("error opening connection with discord: %v", err)
		return
	}

	dg.AddHandler(discord.Ready)
	dg.AddHandler(discord.GuildCreate)
	dg.AddHandler(discord.GuildDelete)
	dg.AddHandler(discord.VoiceStateUpdate)
	dg.AddHandler(discord.ChannelCreate)
	dg.AddHandler(discord.GuildMemberUpdate)

	// Register commands
	discord.RegisterCommands(dg, commands.CommandHandlers)

	// Register interaction handler
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User == nil {
			if _, ok := BlackList[i.Member.User.ID]; ok {
				embed.SendAndDeleteEmbedInteraction(
					s,
					embed.NewEmbed().SetTitle(s.State.User.Username).
						AddField("Error", "You are blacklisted from using this bot", false).
						SetColor("red").
						MessageEmbed, i.Interaction, 5*time.Second)
			} else {
				if command, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
					command.Handler(s, i)
				}
			}
		}
	})

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	log.Infof("%s#%s is now running. Press Ctrl+C to exit", dg.State.User.Username, dg.State.User.Discriminator)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
