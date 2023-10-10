package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	CMD "github.com/nabomhalang/halangcordgo/cmd"
	"github.com/nabomhalang/halangcordgo/config"
	"github.com/nabomhalang/halangcordgo/embed"
	"github.com/nabomhalang/halangcordgo/queue"
	"github.com/nabomhalang/halangcordgo/sponsorblock"
	"github.com/nabomhalang/halangcordgo/youtube"
)

var (
	Youtube *youtube.Youtube = youtube.YT
	yt      youtube.YtDLP
	log     *config.Logger = config.NewLogger("utils")
)

func SearchSong(song string) (*youtube.Video, error) {
	if Youtube != nil {
		result := Youtube.Search(song, 1)
		if len(result) > 0 {
			return &youtube.Video{
				ID:          result[0].ID,
				Title:       result[0].Title,
				Description: result[0].Description,
				Thumbnail:   result[0].Thumbnail,
				Channel:     result[0].Channel,
				Duration:    result[0].Duration,
				Link:        result[0].Link,
			}, nil
		}
	}

	return &youtube.Video{}, errors.New("No results found")
}

func DownloadAndPlay(s *discordgo.Session, guildID, username string, video *youtube.Video, i *discordgo.Interaction, random, loop, respond, priority bool) {
	var c chan struct{}
	if respond {
		c = make(chan struct{})
		go embed.SendEmbedInteraction(s,
			embed.NewEmbed().
				SetTitle(fmt.Sprintf("Add song - %s", video.Title)).
				SetImage(video.Thumbnail).
				SetURL(video.Link).
				SetAuthor(fmt.Sprintf("Request by %s", i.Member.User.Username), i.Member.User.AvatarURL("2048"), "").
				AddField("Channel", video.Channel, true).
				AddField("Duration", SecondsToHHMMSS(int(video.Duration)), true).
				SetColor("green").
				MessageEmbed, i, c)
	}

	link := cleanURL(video.Link)

	info, err := GetInfo(link)
	if err != nil {
		embed.SendAndDeleteEmbedInteraction(s,
			embed.NewEmbed().
				SetTitle("Error").
				AddField("error", err.Error(), false).
				SetColor("red").
				MessageEmbed, i, time.Second*10)
		return
	}

	if random {
		info = Shuffle(info)
	}

	if respond {
		go embed.DeleteInteraction(s, i, c)
	}

	elements := make([]queue.Elements, 0, len(info))
	for _, sj := range info {
		_ = json.Unmarshal([]byte(sj), &yt)

		el := queue.Elements{
			Title:       yt.Title,
			Duration:    SecondsToHHMMSS(int(yt.Duration)),
			Link:        yt.WebpageURL,
			User:        username,
			Thumbnail:   yt.Thumbnail,
			TextChannel: i.ChannelID,
			Loop:        loop,
		}

		// exist := false
		switch yt.Extractor {
		case "youtube":
			el.ID = fmt.Sprintf("%s-%s", yt.ID, yt.Extractor)
			el.Segments = sponsorblock.GetSegments(yt.ID)

			// exist = true
			el.Link = "https://youtu.be/" + yt.ID
		case "generic":
			el.ID = IdGen(el.Link) + "-" + yt.Extractor
		default:
			el.ID = yt.ID + "-" + yt.Extractor
		}

		cacheInfo, err := os.Stat(fmt.Sprintf("%s/%s.dca", config.Get().CachePath, el.ID))
		if err != nil || cacheInfo.Size() <= 0 {
			pipe, cmd := CMD.Gen(yt.WebpageURL, el.ID, CheckAudioOnly(yt.RequestedFormats))
			el.Reader = pipe
			el.Downloading = true

			el.BeforePlay = func() {
				CMD.CmdsStart(cmd)
			}

			el.AfterPlay = func() {
				CMD.CmdsWait(cmd)
			}
		} else {
			f, _ := os.Open(fmt.Sprintf("%s%s.dca", config.Get().CachePath, el.ID))
			el.Reader = f
			el.Closer = f
		}

		elements = append(elements, el)
	}

	Server[guildID].AddSong(s, priority, elements...)
}
