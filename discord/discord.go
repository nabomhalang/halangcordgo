package discord

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/handler"
)

var (
	log *config.Logger = config.NewLogger("discord")
)

func Init() {
	log.Info("Initializing Discord...")

	dg, err := discordgo.New("Bot " + config.GetEnv().TOKEN)
	if err != nil {
		log.Fatalf("error creating Discord session, %v", err)
		return
	}
	defer dg.Close()

	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection, %v", err)
		return
	}

	dg.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name:  "모든 사람들",
				Type:  discordgo.ActivityTypeCustom,
				State: fmt.Sprintf(`"%s %s"을 입력하세요 `, config.GetEnv().Prefix, "도움"),
			},
		},
		Status: "online",
	})

	Register(dg, handler.CommandHandlers)

	log.Infof("%s#%s is now running. Press Ctrl+C to exit", dg.State.User.Username, dg.State.User.Discriminator)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
