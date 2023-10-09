package server

import (
	"sync"
	"sync/atomic"

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

func (s *Server) AddSong(priority bool, el ...queue.Elements) {
	if priority {
		s.Queue.AddPriority(el...)
	} else {
		s.Queue.Add(el...)
	}

	if s.Started.CompareAndSwap(false, true) {
		go s.play()
	}
}

func (s *Server) play() {
	s.Paused.Store(false)

	for el := s.Queue.GetFirst(); el != nil && !s.Clear.Load(); el = s.Queue.GetFirst() {
		go func() {
		}()
	}
}

func (s *Server) IsPlaying() bool {
	return s.Started.Load() && !s.Queue.IsEmpty()
}

func (s *Server) ClearQueue() {
	if s.IsPlaying() {
		s.Clear.Store(true)
		s.Skip <- Clear

		s.Wg.Wait()
		s.Clear.Store(false)

		q := s.Queue.GetAll()
		s.Queue.Clear()

		for _, el := range q {
			if el.Closer != nil {
				el.Closer.Close()
			}
		}
	}
}
