package queue

import (
	"io"
	"sync"
)

type Elements struct {
	ID          string
	Title       string
	Duration    string
	Link        string
	User        string
	Thumbnail   string
	Segments    map[int]bool
	Reader      io.Reader
	Closer      io.Closer
	Downloading bool
	TextChannel string
	BeforePlay  func()
	AfterPlay   func()
	Loop        bool
}

type Queue struct {
	queue []Elements
	rw    *sync.RWMutex
}
