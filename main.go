package main

import (
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/discord"
)

func main() {
	config.Init()

	logger := config.NewLogger("main")

	discord.Init()

	logger.Fatalf("exiting...")
}
