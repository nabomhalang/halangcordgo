package utils

import (
	"os/exec"

	"github.com/nabomhalang/halangcordgo/server"
)

var (
	Server = server.SV
)

func IsCommandNotAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err != nil
}

func ClearAndExit(guildID string) {
	Server[guildID].ClearQueue()

	if Server[guildID].VC != nil {
		Server[guildID].VC.Disconnect()
		Server[guildID].VC = nil
		Server[guildID].VoiceChannel = ""
	}
}

func InitServer(guildID string) {
	if _, ok := Server[guildID]; !ok {
		Server[guildID] = server.NewServer(guildID)
	}
}
