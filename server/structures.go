package server

import (
	"sync"
	"sync/atomic"

	"github.com/bwmarrin/discordgo"
	"github.com/nabomhalang/halangcordgo/queue"
)

type Server struct {
	VC                  *discordgo.VoiceConnection
	VoiceChannel        string
	VoiceChannelMembers map[string]*atomic.Int32
	Started             atomic.Bool
	Clear               atomic.Bool
	Paused              atomic.Bool
	GuildID             string
	Wg                  *sync.WaitGroup
	Pause               chan struct{}
	Resume              chan struct{}
	Frames              int
	Queue               queue.Queue
	Skip                chan SkipReason
}

type SkipReason int

const (
	Error SkipReason = iota
	Finished
	Skip
	Clear
)
