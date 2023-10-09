package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
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
	yt        *youtube.Youtube
	log       *config.Logger = config.NewLogger("main")
)

func init() {
	var err error
	log.Info("Initializing HalangcordGo...")

	if _, err := os.Stat(config.Get().CachePath); err != nil {
		if err = os.Mkdir(config.Get().CachePath, 0755); err != nil {
			log.Errorf("Failed to create cache directory: %s", err.Error())
		}
	}

	os.Remove("--Frag1")

	if utils.IsCommandNotAvailable("ffmpeg") {
		log.Fatal("ffmpeg is not installed")
	}

	if utils.IsCommandNotAvailable("yt-dlp") {
		log.Fatal("yt-dlp is not installed")
	}

	if config.Get().YoutubeKey != "" {
		yt, err = youtube.NewYoutube(config.Get().YoutubeKey)
		if err != nil {
			log.Fatalf("Failed to initialize youtube: %s", err.Error())
		}
	}
}

func main() {
	log.Info("Starting HalangcordGo...")

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("Failed to create discord session: %s", err.Error())
		return
	}

	dg.AddHandler(discord.Ready)
	dg.AddHandler(discord.GuildCreate)
	dg.AddHandler(discord.GuildDelete)
	dg.AddHandler(discord.VoiceStateUpdate)
	dg.AddHandler(discord.ChannelCreate)
	dg.AddHandler(discord.GuildMemberUpdate)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.User == nil {
			if _, ok := BlackList[i.Member.User.ID]; ok {
				embed.SendAndDeleteEmbedInteraction(
					s,
					embed.NewEmbed().SetTitle(s.State.User.Username).AddField("Error",
						"You are blacklisted from using this bot").SetColor(0xff0000).MessageEmbed, i.Interaction, 5*time.Second)
			} else {
				// if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				// 	h(s, i)
				// }
			}
		}
	})

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Failed to open discord session: %s", err.Error())
		return
	}

	s = dg

	log.Info("HalangcordGo is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Error("HalangcordGo is now shutting down...")
	_ = dg.Close()
}
