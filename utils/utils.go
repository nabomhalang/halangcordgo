package utils

import (
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/embed"
	"github.com/nabomhalang/halangcordgo/server"
	"github.com/nabomhalang/halangcordgo/youtube"
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

func UserIsInVoiceChannel(s *discordgo.Session, guildID, userID string) *discordgo.VoiceState {
	g, err := s.State.Guild(guildID)
	if err == nil {
		for _, vs := range g.VoiceStates {
			if vs.UserID == userID {
				return vs
			}
		}
	}

	return nil
}

func JoinVoiceChannel(s *discordgo.Session, i *discordgo.Interaction, channelID string) bool {
	if Server[i.GuildID].VC == nil {
		// Join the voice channel
		vc, err := s.ChannelVoiceJoin(i.GuildID, channelID, false, true)
		if err != nil {
			embed.SendAndDeleteEmbedInteraction(s,
				embed.NewEmbed().
					SetTitle(s.State.User.Username).
					AddField("ERROR", "cannot join voice channel", false).
					SetColor("red").MessageEmbed, i, time.Second*5)
			return false
		}
		Server[i.GuildID].VC = vc
		Server[i.GuildID].VoiceChannel = channelID
	}
	return true
}

func IsValidURL(ul string) bool {
	_, err := url.ParseRequestURI(ul)
	return err == nil
}

func SecondsToHHMMSS(totalSeconds int) string {
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func cleanURL(link string) string {
	u, _ := url.Parse(link)
	q := u.Query()

	q.Del("utm_source")
	q.Del("feature")

	u.RawQuery = q.Encode()

	return u.String()
}

func GetInfo(link string) ([]string, error) {
	// Gets info about songs
	out, err := exec.Command("yt-dlp", "--ignore-errors", "-q", "--no-warnings", "-j", link).CombinedOutput()

	// Parse output as string, splitting it on every newline
	splittedOut := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

	if err != nil {
		return nil, errors.New("Can't get info about song: " + splittedOut[len(splittedOut)-1])
	}

	// Check if yt-dlp returned something
	if strings.TrimSpace(splittedOut[0]) == "" {
		return nil, errors.New("yt-dlp returned no songs")
	}

	return splittedOut, nil
}

func Shuffle(a []string) []string {
	final := make([]string, len(a))

	for i, v := range rand.Perm(len(a)) {
		final[v] = a[i]
	}
	return final
}

func IdGen(link string) string {
	h := sha1.New()
	h.Write([]byte(link))

	return strings.ToLower(base32.HexEncoding.EncodeToString(h.Sum(nil))[0:11])
}

func CheckAudioOnly(formats youtube.RequestedFormats) bool {
	for _, f := range formats {
		if f.Resolution == "audio only" {
			return true
		}
	}

	return false
}
