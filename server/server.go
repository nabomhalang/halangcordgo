package server

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/embed"
	"github.com/nabomhalang/halangcordgo/queue"
)

var (
	SV = make(map[string]*Server)
)

func NewServer(guildID string) *Server {
	return &Server{
		GuildID:             guildID,
		VoiceChannelMembers: make(map[string]*atomic.Int32),
		Wg:                  &sync.WaitGroup{},
		Pause:               make(chan struct{}),
		Resume:              make(chan struct{}),
		Started:             atomic.Bool{},
		Clear:               atomic.Bool{},
		Paused:              atomic.Bool{},
		Queue:               queue.NewQueue(),
	}
}

func (m *Server) AddSong(s *discordgo.Session, priority bool, el ...queue.Elements) {
	if priority {
		m.Queue.AddPriority(el...)
	} else {
		m.Queue.Add(el...)
	}

	if m.Started.CompareAndSwap(false, true) {
		go m.Play(s)
	}
}

func (m *Server) Play(s *discordgo.Session) {
	msg := make(chan *discordgo.Message)
	m.Paused.Store(false)

	for el := m.Queue.GetFirst(); el != nil && !m.Clear.Load(); el = m.Queue.GetFirst() {
		go func() {
			msg <- embed.SendEmbed(s, embed.NewEmbed().SetTitle(s.State.User.Username).
				AddField("Now playing",
					fmt.Sprintf("[%s](%s) - %s added by %s", el.Title, el.Link, el.Duration, el.User), false).
				SetColor(0x7289DA).SetThumbnail(el.Thumbnail).MessageEmbed, el.TextChannel)
		}()

		if el.BeforePlay != nil {
			el.BeforePlay()
		}

		skipReason, _ := PlaySound(m.GuildID, el)

		if el.Downloading && skipReason > Finished {
			devnull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0755)
			_, _ = io.Copy(devnull, el.Reader)
			_ = devnull.Close()
		}

		if el.AfterPlay != nil {
			el.AfterPlay()
		}

		go func() {
			if message := <-msg; message != nil {
				_ = s.ChannelMessageDelete(message.ChannelID, message.ID)
			}
		}()

		if skipReason != Clear {
			m.Queue.RemoveFirst()
		}
	}

	m.Started.Store(false)

	go quitVC(m.GuildID)
}

func (m *Server) IsPlaying() bool {
	return m.Started.Load() && !m.Queue.IsEmpty()
}

func (m *Server) ClearQueue() {
	if m.IsPlaying() {
		m.Clear.Store(true)
		m.Skip <- Clear

		m.Wg.Wait()
		m.Clear.Store(false)

		q := m.Queue.GetAll()
		m.Queue.Clear()

		for _, el := range q {
			if el.Closer != nil {
				el.Closer.Close()
			}
		}
	}
}

func quitVC(guildID string) {
	time.Sleep(1 * time.Minute)

	if SV[guildID].Queue.IsEmpty() && SV[guildID].VC != nil {
		_ = SV[guildID].VC.Disconnect()
		SV[guildID].VC = nil
		SV[guildID].VoiceChannel = ""
	}
}
